package worker

import (
	"context"
	"encoding/json"
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

func New(p usecase.Repository, j usecase.JobRepository, l *zap.Logger, c *config.Config) *Worker {
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
				err := w.makeRequest(job)
				if err != nil {
					w.logger.Error("error accrual: " + err.Error())
					continue
				}
			case <-w.done:
				break loop
			}
		}
	}()

	scheduler := w.scheduler()

	<-w.done

	scheduler.Stop()
}

func (w *Worker) scheduler() *time.Ticker {
	ticker := time.NewTicker(time.Second * 5)

	go func() {
		for {
			select {
			case <-ticker.C:
				err := w.runJob()
				if err != nil {
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
	if err != nil {
		return ErrNoAccessToAccrual
	}

	switch response.StatusCode {
	case http.StatusTooManyRequests:
		return ErrAccrualOverloaded
	case http.StatusInternalServerError:
		return ErrNoAccessToAccrual
	case http.StatusNotFound:
		return ErrOrderNotFound
	case http.StatusNoContent:
		return ErrOrderNotFound
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var result accrualResponse

	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if result.Status == "REGISTERED" || result.Status == "PROCESSING" {
		err = w.updateOrderStatus(result)
		if err != nil {
			return err
		}
		return ErrOrderIsInProcessing
	}

	if result.Status == "INVALID" {
		err = w.updateOrderStatus(result)
		if err != nil {
			return err
		}
		_ = w.closeJob(job)
		return ErrOrderIsInvalid
	}

	err = w.updateOrderStatus(result)
	if err != nil {
		return err
	}

	err = w.closeJob(job)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) updateOrderStatus(result accrualResponse) error {
	w.logger.Info(result.Order + ":" + result.Status)

	order := entity.Order{
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
