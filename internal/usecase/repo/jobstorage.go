package repo

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/nickklius/go-loyalty/internal/entity"
	"github.com/nickklius/go-loyalty/internal/storage/jobstorage"
	"github.com/nickklius/go-loyalty/internal/storage/postgres"
	"github.com/nickklius/go-loyalty/internal/usecase"
)

type JobRepository struct {
	*jobstorage.JobStorage
	*postgres.Postgres
}

func NewJobRepository(j *jobstorage.JobStorage, pg *postgres.Postgres) *JobRepository {
	return &JobRepository{j, pg}
}

func (j *JobRepository) AddJob(_ context.Context, job entity.Job) error {
	j.Lock()
	defer j.Unlock()

	if _, ok := j.Queue[job.OrderID]; ok {
		return usecase.ErrDBDuplicatedEntry
	}

	j.Queue[job.OrderID] = job

	return nil
}

func (j *JobRepository) syncJobQueue(ctx context.Context) error {
	sql, args, err := j.Builder.
		Select("number, status").
		From("orders").
		Where(squirrel.Eq{"status": []string{
			string(entity.OrderStatusNew),
			string(entity.OrderStatusRegistered),
			string(entity.OrderStatusProcessing)}}).
		ToSql()
	if err != nil {
		return fmt.Errorf("jobrepo - getOrdersForJobQueue - r.Builder: %w", err)
	}

	rows, err := j.Pool.Query(ctx, sql, args...)
	if err != nil {
		return err
	}

	for rows.Next() {
		var order entity.Order

		err = rows.Scan(&order.Number, &order.OrderStatus)
		if err != nil {
			return err
		}

		if _, ok := j.Queue[order.Number]; !ok {
			job := entity.Job{
				OrderID: order.Number,
				Status:  order.OrderStatus,
			}

			j.Queue[job.OrderID] = job
		}
	}

	return nil
}

func (j *JobRepository) GetJobs(ctx context.Context) ([]entity.Job, error) {
	j.Lock()
	defer j.Unlock()

	var result []entity.Job

	err := j.syncJobQueue(ctx)
	if err != nil {
		return result, err
	}

	for _, job := range j.Queue {
		result = append(result, job)
	}

	return result, nil
}

func (j *JobRepository) DeleteJob(job entity.Job) error {
	j.Lock()
	defer j.Unlock()

	delete(j.Queue, job.OrderID)
	return nil
}
