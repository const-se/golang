package handler

import "encoding/json"

type APIShortURLResponse struct {
	Result string `json:"result"`
}

func (h *handler) apiShortURLResponse(dto APIShortURLResponse) (string, error) {
	response, err := json.Marshal(dto)
	if err != nil {
		return "", err
	}

	return string(response), nil
}
