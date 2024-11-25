package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/kiquetal/orchestration-go-scratch/task"
	"github.com/kiquetal/orchestration-go-scratch/worker"
	"log"
	"net/http"
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
		selectedWorker := m.SelectWorker()
		e := m.Pending.Dequeue().(*task.TaskEvent)
		t := e.Task
		log.Printf("Sending task %s to selectedWorker %s", t.ID, selectedWorker)
		m.EventDb[e.ID] = e
		m.WorkerTasks[selectedWorker] = append(m.WorkerTasks[selectedWorker], e.ID)
		m.TaskWorkerMap[e.ID] = selectedWorker
		t.State = task.Scheduled
		m.TaskDb[t.ID] = &t
		data, err := json.Marshal(e)
		if err != nil {
			log.Printf("Error marshalling task event: %s", err)
		}

		url := fmt.Sprintf("https://%s/tasks", selectedWorker)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Printf("Error sending task to selectedWorker: %s", err)
			m.Pending.Enqueue(e)
			return

		}
		d := json.NewDecoder(resp.Body)
		if resp.StatusCode != http.StatusCreated {
			e := worker.ErrResponse{}
			err := d.Decode(&e)
			if err != nil {
				log.Printf("Error decoding error response: %s", err)
				return
			}
			log.Printf("Error response from selectedWorker: %s", e.HTTPStatusCode, e.Message)
			return
		}
		t = task.Task{}
		err = d.Decode(&t)
		if err != nil {
			log.Printf("Error decoding task response: %s", err)
			return
		}
		log.Print("%#v\n", t)
	} else {
		log.Println("No tasks to send")
	}
}

func (m *Manager) UpdateTasks() {
	for _, w := range m.Workers {
		log.Printf("Checking w %v for tasks", w)
		url := fmt.Sprintf("https://%s/tasks", w)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error getting tasks from w: %s", err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("Error getting tasks from w: %d", resp.StatusCode)

		}
		d := json.NewDecoder(resp.Body)
		var tasks []task.Task
		err = d.Decode(&tasks)
		if err != nil {
			log.Printf("Error decoding tasks from w: %s", err)
		}

		for _, t := range tasks {
			log.Printf("Attempt to update task %s", t.ID)
			taskInDb, ok := m.TaskDb[t.ID]
			if !ok {
				log.Printf("Task %s not found in db", t.ID)
				return
			}
			if taskInDb.State != t.State {
				log.Printf("Task %s state changed from %s to %s", t.ID, taskInDb.State, t.State)
				m.TaskDb[t.ID].State = t.State
			}

			m.TaskDb[t.ID].StartTime = t.StartTime
			m.TaskDb[t.ID].FinishTime = t.FinishTime
			m.TaskDb[t.ID].ContainerID = t.ContainerID

		}

	}
}
