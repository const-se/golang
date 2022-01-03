package handler

import (
	"net/http"
	"strconv"
)

func (h *handler) idFromRequest(request *http.Request) (int, error) {
	id, err := strconv.Atoi(request.URL.Path[1:])
	if err != nil {
		return 0, err
	}

	return id, nil
}
