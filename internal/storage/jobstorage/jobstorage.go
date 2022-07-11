package jobstorage

import (
	"sync"

	"github.com/nickklius/go-loyalty/internal/entity"
)

type JobStorage struct {
	sync.Mutex
	Queue map[string]entity.Job
}

func NewJobStorage() *JobStorage {
	return &JobStorage{
		Queue: map[string]entity.Job{},
	}
}
