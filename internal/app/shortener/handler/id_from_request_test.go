package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handler_idFromRequest(t *testing.T) {
	type args struct {
		request *http.Request
	}

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Корректный запрос",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/123", nil),
			},
			want: 123,
		},
		{
			name: "Пустой запрос",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			wantErr: true,
		},
		{
			name: "Некорректный запрос",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/abc", nil),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler{}

			got, err := h.idFromRequest(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("idFromRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("idFromRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
