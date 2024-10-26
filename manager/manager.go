package manager

import (
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/kiquetal/orchestration-go-scratch/task"
)

type Manager struct {
	Pending       *queue.Queue
	TaskDb        map[string][]*task.Task
	EventDb       map[string][]*task.TaskEvent
	Workers       []string
	WorkerTasks   map[string][]uuid.UUID
	TaskWorkerMap map[uuid.UUID]string
}

func (m *Manager) SelectWorker() string {
	// Select a worker based on some strategy
	return m.Workers[0]
}

func (m *Manager) UpdateTask() {
	// Update task state
	fmt.Println("Task Updated")
}

func (m *Manager) SendWork() {
	fmt.Println("Work Sent")
}
