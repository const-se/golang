package handler

import (
	"github.com/const-se/golang/internal/app/shortener/repository"
	"github.com/go-chi/chi/v5"
)

const (
	ContentTypeHeader = "Content-Type"
	locationHeader    = "Location"

	ContentTypeJSON = "application/json"
)

type handler struct {
	*chi.Mux
	repository repository.Repository
	baseURL    string
}

func NewHandler(repository repository.Repository, baseURL string) *handler {
	h := &handler{
		Mux:        chi.NewMux(),
		repository: repository,
		baseURL:    baseURL,
	}

	h.Post("/", h.shortURL())
	h.Post("/api/shorten", h.apiShortURL())
	h.Get("/{id:\\d+}", h.unshortURL())

	return h
}
