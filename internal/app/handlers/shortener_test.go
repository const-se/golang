package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortener_Handle(t *testing.T) {
	handler := new(Shortener)

	type want struct {
		statusCode   int
		headers      map[string][]string
		responseBody string
	}

	type args struct {
		method string
		path   string
		body   io.Reader
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Создание короткой ссылки",
			args: args{
				method: http.MethodPost,
				path:   "/",
				body:   strings.NewReader("https://practicum.yandex.ru"),
			},
			want: want{
				statusCode:   http.StatusCreated,
				responseBody: "http://localhost:8080/0",
			},
		},
		{
			name: "Получение исходной ссылки",
			args: args{
				method: http.MethodGet,
				path:   "/0",
			},
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				headers: map[string][]string{
					"Location": {
						"https://practicum.yandex.ru",
					},
				},
			},
		},
		{
			name: "Некорректный запрос",
			args: args{
				method: http.MethodGet,
				path:   "/",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.args.method, tt.args.path, tt.args.body)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.Handle)
			h.ServeHTTP(w, r)
			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)

			for header, values := range tt.want.headers {
				assert.Equal(t, values, res.Header.Values(header))
			}

			responseBody, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)
			err = res.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.responseBody, string(responseBody))
		})
	}
}
