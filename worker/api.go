package worker

import (
	"github.com/go-chi/chi/v5"
)

type Api struct {
	Address string
	Port    int
	Worker  *Worker
	Router  *chi.Mux
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
}
