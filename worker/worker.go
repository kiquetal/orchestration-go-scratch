package worker

import (
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/kiquetal/orchestration-go-scratch/task"
)

type Worker struct {
	Name      string
	Queue     *queue.Queue
	Db        map[uuid.UUID]*task.Task
	TaskCount int
}

func (w *Worker) CollectStats() {
	fmt.Println("Worker: ", w.Name)
	fmt.Println("Task Count: ", w.TaskCount)
	fmt.Println("Queue Length: ", w.Queue.Len())
}

func (w *Worker) RunTask(t *task.Task) {
	w.Queue.Enqueue(t)
	w.TaskCount++
}

func (w *Worker) StartTask() {
	if w.Queue.Len() > 0 {
		t := w.Queue.Dequeue().(*task.Task)
		w.TaskCount--
		w.Db[t.ID] = t
		fmt.Println("Task Started: ", t.Name)
	}
}

func (w *Worker) StopTask(t *task.Task) {
	delete(w.Db, t.ID)
	fmt.Println("Task Stopped: ", t.Name)
}
