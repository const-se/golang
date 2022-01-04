package handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_handler_apiShortURLResponse(t *testing.T) {
	type args struct {
		dto APIShortURLResponse
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Корректный ответ",
			args: args{
				dto: APIShortURLResponse{
					Result: "https://practicum.yandex.ru",
				},
			},
			want:    "{\"result\":\"https://practicum.yandex.ru\"}",
			wantErr: assert.NoError,
		},
		{
			name:    "Пустой ответ",
			want:    "{\"result\":\"\"}",
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler{}

			got, err := h.apiShortURLResponse(tt.args.dto)
			if !tt.wantErr(t, err, fmt.Sprintf("apiShortURLResponse(%v)", tt.args.dto)) {
				return
			}

			assert.Equalf(t, tt.want, got, "apiShortURLResponse(%v)", tt.args.dto)
		})
	}
}
