package scheduler

import (
	"fmt"
	"github.com/kiquetal/orchestration-go-scratch/worker"
)

type Scheduler interface {
	SelectCandidateNodes()
	Score()
	Pick()
}

type SchedulerImpl struct {
	Workers []worker.Worker
}

func (s *SchedulerImpl) AddWorker(worker worker.Worker) {
	fmt.Println("Worker Added: ", worker.Name)
	s.Workers = append(s.Workers, worker)

}
