package handler

import (
	"fmt"
	"net/http"
)

func (h *handler) apiShortURL() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		apiURL, err := h.apiShortURLRequest(request)
		if err != nil {
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(apiURL.URL) == 0 {
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := h.repository.SaveURL(apiURL.URL)
		if err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			return
		}

		response, err := h.apiShortURLResponse(APIShortURLResponse{
			Result: fmt.Sprintf("%s/%d", baseURL, id),
		})
		if err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			return
		}

		responseWriter.Header().Set(ContentTypeHeader, ContentTypeJSON)
		responseWriter.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprint(responseWriter, response)
	}
}
