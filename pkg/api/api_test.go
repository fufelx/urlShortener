package api

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddUrl(t *testing.T) {
	InMemory = true
	UrlToShorturl = make(map[string]string)
	ShorturlToUrl = make(map[string]string)

	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Ошибка: неверный метод",
			method:     http.MethodGet,
			body:       "",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "неправильный метод\n",
		},
		{
			name:       "Ошибка: некорректный JSON",
			method:     http.MethodPost,
			body:       `{"wrong_field": "https://example.com"}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   "неправильный JSON\n",
		},
		{
			name:       "Успешное сокращение ссылки",
			method:     http.MethodPost,
			body:       `{"url": "https://example.com"}`,
			wantStatus: http.StatusOK,
			wantBody:   `"short123"\n`,
		},
		{
			name:       "Ошибка БД",
			method:     http.MethodPost,
			body:       `{"url": "https://error.com"}`,
			wantStatus: http.StatusInternalServerError,
			wantBody:   "ошибка при сокращении ссылки\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/add", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			AddUrl(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode)

			body := new(bytes.Buffer)
			body.ReadFrom(res.Body)
			assert.Equal(t, tt.wantBody, body.String())
		})
	}
}

func TestGetUrl(t *testing.T) {
	InMemory = false

	tests := []struct {
		name       string
		shorturl   string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Ошибка: пустой параметр",
			shorturl:   "",
			wantStatus: http.StatusBadRequest,
			wantBody:   "shorturl не найдена в RawQuery\n",
		},
		{
			name:       "Успешное получение оригинальной ссылки",
			shorturl:   "short123",
			wantStatus: http.StatusOK,
			wantBody:   `"https://example.com"\n`,
		},
		{
			name:       "Ошибка: ссылка не найдена",
			shorturl:   "notfound",
			wantStatus: http.StatusBadRequest,
			wantBody:   "оригинал ссылки отсутствует\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/get?shorturl="+tt.shorturl, nil)
			rec := httptest.NewRecorder()

			GetUrl(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode)

			body := new(bytes.Buffer)
			body.ReadFrom(res.Body)
			assert.Equal(t, tt.wantBody, body.String())
		})
	}
}
