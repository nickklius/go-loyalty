package worker

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/nickklius/go-loyalty/config"
	"github.com/nickklius/go-loyalty/internal/entity"
	"github.com/nickklius/go-loyalty/internal/usecase"
)

type Worker struct {
	pg     usecase.Repository
	repo   usecase.JobRepository
	logger *zap.Logger
	cfg    *config.Config
	stream chan entity.Job
	done   chan struct{}
}

type accrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func NewWorker(p usecase.Repository, j usecase.JobRepository, l *zap.Logger, c *config.Config) *Worker {
	return &Worker{
		pg:     p,
		repo:   j,
		logger: l,
		cfg:    c,
		stream: make(chan entity.Job),
		done:   make(chan struct{}),
	}
}

func (w *Worker) Run() {
	go func() {
	loop:
		for {
			select {
			case job := <-w.stream:
				w.logger.Info("job is come")
				err := w.makeRequest(job)
				if err != nil {
					w.logger.Error("error when making request to the accrual " + err.Error())
					continue
				}
			case <-w.done:
				break loop
			}
		}
	}()

	sched := w.scheduler()

	<-w.done

	sched.Stop()
}

func (w *Worker) scheduler() *time.Ticker {
	ticker := time.NewTicker(time.Second * 5)
	w.logger.Info("start scheduler")

	go func() {
		for {
			select {
			case <-ticker.C:
				err := w.runJob()
				if err != nil {
					w.logger.Error("error when job running" + err.Error())
					continue
				}
			case <-w.done:
				return
			}
		}
	}()

	return ticker
}

func (w *Worker) Done() {
	w.done <- struct{}{}
}

func (w *Worker) pushJob(job entity.Job) {
	w.stream <- job
}

func (w *Worker) runJob() error {
	ticker := time.NewTicker(time.Second * 1)
	jobs, err := w.repo.GetJobs()
	if err != nil {
		return err
	}

	for _, j := range jobs {
		w.pushJob(j)
		<-ticker.C
	}

	ticker.Stop()
	return nil
}

func (w *Worker) closeJob(job entity.Job) error {
	err := w.repo.DeleteJob(job)
	if err != nil {
		return err
	}
	return nil
}

func (w *Worker) makeRequest(job entity.Job) error {
	response, err := http.Get(w.cfg.Accrual.AccrualAddress + "/api/orders/" + job.OrderID)
	defer response.Body.Close()

	if err != nil {
		w.logger.Info("Problem with access accrual service")
		return errors.New("problem with access accrual service")
	}
	if response.StatusCode == http.StatusTooManyRequests {
		w.logger.Info("Accrual service overloaded")
		return errors.New("accrual service overloaded")
	}
	if response.StatusCode == http.StatusInternalServerError {
		w.logger.Info("Accrual service is unavailable")
		return errors.New("accrual service is unavailable")
	}
	if response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusNoContent {
		w.logger.Info("Order not found on accrual service")
		return errors.New("order not found on accrual service")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var result accrualResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if result.Status == "REGISTERED" || result.Status == "PROCESSING" {
		w.logger.Info("order processing is not finished")
		return errors.New("order processing is not finished")
	}

	if result.Status == "INVALID" {
		w.logger.Info("order will not be processed")
		err = w.closeJob(job)
		if err != nil {
			return errors.New("error when removing invalid job from queue")
		}
		return errors.New("order processing is not finished")
	}

	err = w.updateOrderStatus(result)
	if err != nil {
		w.logger.Info("order processing is not finished")
		return errors.New("error in storing order status update")
	}

	err = w.closeJob(job)
	if err != nil {
		return errors.New("error when removing success job from queue")
	}

	return nil
}

func (w *Worker) updateOrderStatus(result accrualResponse) error {
	var order entity.Order

	order = entity.Order{
		Number:  result.Order,
		Status:  result.Status,
		Accrual: result.Accrual,
	}

	err := w.pg.UpdateOrderStatus(context.TODO(), order)
	if err != nil {
		return err
	}
	return nil
}
