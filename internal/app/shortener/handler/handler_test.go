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
