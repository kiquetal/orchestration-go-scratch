package manager

import (
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/kiquetal/orchestration-go-scratch/task"
	"log"
)

type Manager struct {
	Pending       *queue.Queue
	TaskDb        map[uuid.UUID]*task.Task
	EventDb       map[uuid.UUID]*task.TaskEvent
	Workers       []string
	WorkerTasks   map[string][]uuid.UUID
	TaskWorkerMap map[uuid.UUID]string
	LastWorker    int
}

func (m *Manager) SelectWorker() string {
	var newWorker int

	if m.LastWorker+1 < len(m.Workers) {
		newWorker = m.LastWorker + 1
		m.LastWorker++
	} else {
		newWorker = 0
		m.LastWorker = 0

	}
	return m.Workers[newWorker]
}

func (m *Manager) SendWork() {
	if m.Pending.Len() > 0 {
		worker := m.SelectWorker()
		e := m.Pending.Dequeue().(*task.TaskEvent)
		t := e.Task
		log.Printf("Sending task %s to worker %s", t.ID, worker)
		m.EventDb[e.ID] = e
		m.WorkerTasks[worker] = append(m.WorkerTasks[worker], e.ID)
		m.TaskWorkerMap[e.ID] = worker
		t.State = task.Scheduled
		m.TaskDb[t.ID] = &t
	}
}

func (m *Manager) UpdateTask() {
	// Update task state
	fmt.Println("Task Updated")
}
