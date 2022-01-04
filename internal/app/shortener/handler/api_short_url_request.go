package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type APIShortURLRequest struct {
	URL string `json:"url"`
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
