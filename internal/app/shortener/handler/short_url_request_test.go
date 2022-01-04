package handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_handler_shortURLRequest(t *testing.T) {
	type args struct {
		request *http.Request
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Корректный запрос",
			args: args{
				request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://practicum.yandex.ru")),
			},
			want:    "https://practicum.yandex.ru",
			wantErr: assert.NoError,
		},
		{
			name: "Пустой запрос",
			args: args{
				request: httptest.NewRequest(http.MethodPost, "/", nil),
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler{}

			got, err := h.shortURLRequest(tt.args.request)
			if !tt.wantErr(t, err, fmt.Sprintf("shortURLRequest(%v)", tt.args.request)) {
				return
			}

			assert.Equalf(t, tt.want, got, "shortURLRequest(%v)", tt.args.request)
		})
	}
}
