package handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handler_unshortURLRequest(t *testing.T) {
	type args struct {
		request *http.Request
	}

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Корректный запрос",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/123", nil),
			},
			want:    123,
			wantErr: assert.NoError,
		},
		{
			name: "Пустой запрос",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			wantErr: assert.Error,
		},
		{
			name: "Некорректный запрос",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/abc", nil),
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler{}

			got, err := h.unshortURLRequest(tt.args.request)
			if !tt.wantErr(t, err, fmt.Sprintf("unshortURLRequest(%v)", tt.args.request)) {
				return
			}

			assert.Equalf(t, tt.want, got, "unshortURLRequest(%v)", tt.args.request)
		})
	}
}
