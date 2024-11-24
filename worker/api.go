package worker

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type Api struct {
	Address string
	Port    int
	Worker  *Worker
	Router  *chi.Mux
}

type ErrResponse struct {
	HTTPStatusCode int    `json:"status"`
	Message        string `json:"message"`
}

func (a *Api) initRouter() {
	a.Router = chi.NewRouter()
	a.Router.Route("/tasks", func(r chi.Router) {
		r.Post("/", a.StartTaskHandler)
		r.Get("/", a.GetTaskHandler)
		r.Route("/{taslID}", func(r1 chi.Router) {
			r1.Delete("/", a.StopTaskHandler)
		})
	})
	a.Router.Route("/stats", func(r chi.Router) {
		r.Get("/", a.GetStatsHandler)
	})
}

func (a *Api) Start() {
	log.Printf("Starting worker %s on %s:%d", a.Worker.Name, a.Address, a.Port)
	a.initRouter()
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", a.Address, a.Port), a.Router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
