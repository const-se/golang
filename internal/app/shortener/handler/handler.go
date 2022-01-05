package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/const-se/golang/internal/app/shortener/repository"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	ContentTypeHeader = "Content-Type"
	locationHeader    = "Location"

	ContentTypeJSON = "application/json"
)

type APIShortURLRequest struct {
	URL string `json:"url"`
}

type APIShortURLResponse struct {
	Result string `json:"result"`
}

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

func (h *handler) shortURL() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		url, err := h.shortURLRequest(request)
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
		_, _ = fmt.Fprintf(responseWriter, "%s/%d", h.baseURL, id)
	}
}

func (h *handler) shortURLRequest(request *http.Request) (string, error) {
	defer func() {
		_ = request.Body.Close()
	}()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

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
			Result: fmt.Sprintf("%s/%d", h.baseURL, id),
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

func (h *handler) apiShortURLRequest(request *http.Request) (dto APIShortURLRequest, err error) {
	defer func() {
		_ = request.Body.Close()
	}()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &dto)

	return
}

func (h *handler) apiShortURLResponse(dto APIShortURLResponse) (string, error) {
	response, err := json.Marshal(dto)
	if err != nil {
		return "", err
	}

	return string(response), nil
}

func (h *handler) unshortURL() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		id, err := h.unshortURLRequest(request)
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

func (h *handler) unshortURLRequest(request *http.Request) (int, error) {
	id, err := strconv.Atoi(request.URL.Path[1:])
	if err != nil {
		return 0, err
	}

	return id, nil
}
