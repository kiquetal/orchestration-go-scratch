package worker

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/kiquetal/orchestration-go-scratch/task"
	"log"
	"net/http"
)

func (a *Api) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Start a task
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	var t task.TaskEvent
	err := d.Decode(&t)
	if err != nil {
		msg := fmt.Sprintf("Failed to decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		e := types.ErrorResponse{
			Message: msg,
		}
		json.NewEncoder(w).Encode(e)
		return
	}
	a.Worker.AddTask(t.Task)
	log.Printf("Task %s added to worker %s", t.Task.ID, a.Worker.Name)
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(t.Task)
	if err != nil {
		log.Printf("Failed to encode task: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (a *Api) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Stop a task
}

func (a *Api) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Get a task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(a.Worker.GetTasks())
	if err != nil {
		log.Printf("Failed to encode tasks: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
