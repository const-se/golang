package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_handler_urlFormRequest(t *testing.T) {
	type args struct {
		request *http.Request
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Корректный запрос",
			args: args{
				request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://practicum.yandex.ru")),
			},
			want: "https://practicum.yandex.ru",
		},
		{
			name: "Пустой запрос",
			args: args{
				request: httptest.NewRequest(http.MethodPost, "/", nil),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler{}

			got, err := h.urlFormRequest(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("urlFormRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("urlFormRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
