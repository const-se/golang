package handler

import (
	"fmt"
	"github.com/const-se/golang/internal/app/shortener/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	testID      = 123
	testURL     = "https://practicum.yandex.ru"
	testBaseURL = "http://localhost:8080"
)

type testRepository struct {
}

func (r *testRepository) URL(id int) (string, error) {
	if id == testID {
		return testURL, nil
	}

	return "", repository.ErrURLNotFound
}

func (r *testRepository) SaveURL(_ string) (int, error) {
	return testID, nil
}

func TestHandler(t *testing.T) {
	h := NewHandler(&testRepository{}, testBaseURL)
	server := httptest.NewServer(h)
	defer server.Close()
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	type args struct {
		method      string
		path        string
		requestBody io.Reader
	}

	tests := []struct {
		name           string
		args           args
		assertResponse func(t *testing.T, response *http.Response, responseBody string)
	}{
		{
			name: "Создание короткой ссылки",
			args: args{
				method:      http.MethodPost,
				path:        "/",
				requestBody: strings.NewReader(testURL),
			},
			assertResponse: func(t *testing.T, response *http.Response, responseBody string) {
				assert.Equal(t, http.StatusCreated, response.StatusCode)
				assert.Equal(t, fmt.Sprintf("%s/%d", testBaseURL, testID), responseBody)
			},
		},
		{
			name: "Пустой запрос на создание короткой ссылки",
			args: args{
				method: http.MethodPost,
				path:   "/",
			},
			assertResponse: func(t *testing.T, response *http.Response, responseBody string) {
				assert.Equal(t, http.StatusBadRequest, response.StatusCode)
				assert.Empty(t, responseBody)
			},
		},
		{
			name: "Некорректный запрос на создание короткой ссылки",
			args: args{
				method: http.MethodPost,
				path:   "/456",
			},
			assertResponse: func(t *testing.T, response *http.Response, responseBody string) {
				assert.Equal(t, http.StatusMethodNotAllowed, response.StatusCode)
				assert.Empty(t, responseBody)
			},
		},
		{
			name: "Создание короткой ссылки (API)",
			args: args{
				method:      http.MethodPost,
				path:        "/api/shorten",
				requestBody: strings.NewReader(fmt.Sprintf("{\"url\":\"%s\"}", testURL)),
			},
			assertResponse: func(t *testing.T, response *http.Response, responseBody string) {
				assert.Equal(t, http.StatusCreated, response.StatusCode)
				assert.Equal(t, fmt.Sprintf("{\"result\":\"%s/%d\"}", testBaseURL, testID), responseBody)
			},
		},
		{
			name: "Пустой запрос на создание короткой ссылки (API)",
			args: args{
				method: http.MethodPost,
				path:   "/api/shorten",
			},
			assertResponse: func(t *testing.T, response *http.Response, responseBody string) {
				assert.Equal(t, http.StatusBadRequest, response.StatusCode)
				assert.Empty(t, responseBody)
			},
		},
		{
			name: "Некорректный запрос на создание короткой ссылки (API)",
			args: args{
				method:      http.MethodPost,
				path:        "/api/shorten",
				requestBody: strings.NewReader("{url}"),
			},
			assertResponse: func(t *testing.T, response *http.Response, responseBody string) {
				assert.Equal(t, http.StatusBadRequest, response.StatusCode)
				assert.Empty(t, responseBody)
			},
		},
		{
			name: "Получение исходной ссылки",
			args: args{
				method: http.MethodGet,
				path:   fmt.Sprintf("/%d", testID),
			},
			assertResponse: func(t *testing.T, response *http.Response, responseBody string) {
				assert.Equal(t, http.StatusTemporaryRedirect, response.StatusCode)
				assert.Equal(t, testURL, response.Header.Get(locationHeader))
				assert.Empty(t, responseBody)
			},
		},
		{
			name: "Получение несуществующей исходной ссылки",
			args: args{
				method: http.MethodGet,
				path:   "/456",
			},
			assertResponse: func(t *testing.T, response *http.Response, responseBody string) {
				assert.Equal(t, http.StatusNotFound, response.StatusCode)
				assert.Empty(t, response.Header.Get(locationHeader))
				assert.Empty(t, responseBody)
			},
		},
		{
			name: "Некорректный запрос на получение исходной ссылки",
			args: args{
				method: http.MethodGet,
				path:   "/",
			},
			assertResponse: func(t *testing.T, response *http.Response, responseBody string) {
				assert.Equal(t, http.StatusMethodNotAllowed, response.StatusCode)
				assert.Empty(t, responseBody)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, responseBody := handleRequest(t, server, client, tt.args.method, tt.args.path, tt.args.requestBody)
			defer func() {
				_ = response.Body.Close()
			}()

			tt.assertResponse(t, response, responseBody)
		})
	}
}

func handleRequest(t *testing.T, server *httptest.Server, client *http.Client, method, path string, requestBody io.Reader) (*http.Response, string) {
	request, err := http.NewRequest(method, server.URL+path, requestBody)
	require.NoError(t, err)

	response, err := client.Do(request)
	require.NoError(t, err)

	responseBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)

	return response, string(responseBody)
}

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
