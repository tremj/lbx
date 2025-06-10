package api

import "github.com/go-chi/chi/v5"

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/configs", ListConfigsHandler)
	r.Get("/configs/{name}", GetConfigHandler)
	r.Post("/configs", SaveConfigHandler)
	r.Delete("/configs", DeleteConfigHandler)

	return r
}
