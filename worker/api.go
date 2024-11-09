package worker

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Api struct {
	Address string
	Port    int
	Worker  *Worker
	Router  *chi.Mux
}

func StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Start a task
}

func StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Stop a task
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Get a task
}
