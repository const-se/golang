package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

type Shortener struct {
	urls  []string
	mutex sync.Mutex
}

func (h *Shortener) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r.URL.Path)
	case http.MethodPost:
		h.handlePost(w, r.URL.Path, r.Body)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *Shortener) handleGet(w http.ResponseWriter, path string) {
	if len(path) <= 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(path[1:])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if id >= len(h.urls) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", h.urls[id])
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Shortener) handlePost(w http.ResponseWriter, path string, body io.ReadCloser) {
	if path != "/" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := readBody(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.mutex.Lock()
	h.urls = append(h.urls, url)
	id := len(h.urls) - 1
	h.mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprintf(w, "http://localhost:8080/%d", id)
}

func readBody(body io.ReadCloser) (string, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(body); err != nil {
		return "", err
	}

	return buf.String(), nil
}
