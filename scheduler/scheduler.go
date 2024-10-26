package scheduler

import (
	"fmt"
	"github.com/kiquetal/orchestration-go-scratch/worker"
)

type Scheduler interface {
	AddWorker(worker worker.Worker)
}

type SchedulerImpl struct {
	Workers []worker.Worker
}

func (s *SchedulerImpl) AddWorker(worker worker.Worker) {
	fmt.Println("Worker Added: ", worker.Name)
	s.Workers = append(s.Workers, worker)

}
