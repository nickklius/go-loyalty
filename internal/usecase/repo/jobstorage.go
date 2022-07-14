package repo

import (
	"context"

	"github.com/nickklius/go-loyalty/internal/entity"
	"github.com/nickklius/go-loyalty/internal/storage/jobstorage"
	"github.com/nickklius/go-loyalty/internal/usecase"
)

type JobRepository struct {
	*jobstorage.JobStorage
}

func NewJobRepository(j *jobstorage.JobStorage) *JobRepository {
	return &JobRepository{j}
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

func (j *JobRepository) GetJobs() ([]entity.Job, error) {
	j.Lock()
	defer j.Unlock()

	var result []entity.Job

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
