package worker

import (
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/kiquetal/orchestration-go-scratch/task"
	"time"
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

}

func (w *Worker) AddTask(t task.Task) {
	w.Queue.Enqueue(&t)

}
func (w *Worker) StartTask(t task.Task) task.DockerResult {
	t.StartTime = time.Now().UTC()
	config := t.NewConfig()
	d := task.NewDocker(config)
	result := d.Run()
	if result.Error != nil {
		fmt.Printf("Error running task: %s\n", result.Error)
		t.State = task.Failed
		w.Db[t.ID] = &t
		return result
	}
	t.ContainerID = result.ContainerId
	t.State = task.Running
	w.Db[t.ID] = &t
	fmt.Printf("Started container %v, for task %v\n", t.ContainerID, t.Name)
	return result
}

func (w *Worker) StopTask(t task.Task) task.DockerResult {
	config := t.NewConfig()
	d := task.NewDocker(config)
	result := d.Stop(t.ContainerID)
	if result.Error != nil {
		fmt.Printf("Error stopping task: %s\n", result.Error)
	}
	t.FinishTime = time.Now().UTC()
	t.State = task.Completed
	w.Db[t.ID] = &t
	fmt.Printf("Stopped and removed container %v, for task %v\n", t.ContainerID, t.Name)

	return result
}
