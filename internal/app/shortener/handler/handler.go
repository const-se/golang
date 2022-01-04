package handler

import (
	"github.com/const-se/golang/internal/app/shortener/repository"
	"github.com/go-chi/chi"
)

const (
	baseURL = "http://localhost:8080"

	ContentTypeHeader = "Content-Type"
	locationHeader    = "Location"

	ContentTypeJSON = "application/json"
)

type handler struct {
	*chi.Mux
	repository repository.Repository
}

func NewHandler(repository repository.Repository) *handler {
	h := &handler{
		Mux:        chi.NewMux(),
		repository: repository,
	}

	h.Post("/", h.shortURL())
	h.Post("/api/shorten", h.apiShortURL())
	h.Get("/{id:\\d+}", h.unshortURL())

	return h
}
