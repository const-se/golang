package handler

import (
	"errors"
	"github.com/const-se/golang/internal/app/shortener/repository"
	"net/http"
)

func (h *handler) unshortURL() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		id, err := h.idFromRequest(request)
		if err != nil {
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}

		url, err := h.repository.URL(id)
		if err != nil {
			if errors.Is(err, repository.ErrURLNotFound) {
				responseWriter.WriteHeader(http.StatusNotFound)
			} else {
				responseWriter.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		responseWriter.Header().Set(locationHeader, url)
		responseWriter.WriteHeader(http.StatusTemporaryRedirect)
	}
}
