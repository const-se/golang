package handler

import (
	"fmt"
	"net/http"
)

func (h *handler) shortURL() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		url, err := h.urlFormRequest(request)
		if err != nil {
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(url) == 0 {
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := h.repository.SaveURL(url)
		if err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			return
		}

		responseWriter.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(responseWriter, "%s/%d", baseURL, id)
	}
}
