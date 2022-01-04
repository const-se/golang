package handler

import (
	"io/ioutil"
	"net/http"
)

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
