package handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_handler_apiShortURLRequest(t *testing.T) {
	type args struct {
		request *http.Request
	}

	tests := []struct {
		name    string
		args    args
		wantDto APIShortURLRequest
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Корректный запрос",
			args: args{
				request: httptest.NewRequest(
					http.MethodPost,
					"/api/shorten",
					strings.NewReader("{\"url\":\"https://practicum.yandex.ru\"}"),
				),
			},
			wantDto: APIShortURLRequest{
				URL: "https://practicum.yandex.ru",
			},
			wantErr: assert.NoError,
		},
		{
			name: "Пустой запрос",
			args: args{
				request: httptest.NewRequest(http.MethodPost, "/api/shorten", nil),
			},
			wantErr: assert.Error,
		},
		{
			name: "Некорректный запрос",
			args: args{
				request: httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader("{url}")),
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler{}

			gotDto, err := h.apiShortURLRequest(tt.args.request)
			if !tt.wantErr(t, err, fmt.Sprintf("apiShortURLRequest(%v)", tt.args.request)) {
				return
			}

			assert.Equalf(t, tt.wantDto, gotDto, "apiShortURLRequest(%v)", tt.args.request)
		})
	}
}
