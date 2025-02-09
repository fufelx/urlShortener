package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddUrl(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		inmemory       bool
		expectedStatus int
	}{
		{"Стандарт", http.MethodPost, `{"url": "https://example.com"}`, true, 200},
		{"Неправильный метод", http.MethodGet, `{"url": "https://example.com"}`, true, 405},
		{"Неправильный JSONx1", http.MethodPost, `{"urlshort": "https://example.com"}`, true, 400},
		{"Неправильный URL", http.MethodPost, `{"url": "example.com"}`, true, 400},
		{"Неправильный JSONx2", http.MethodPost, `{"url": "https://example.com"`, true, 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.inmemory {
				InMemory = true
			}

			body := bytes.NewBufferString(tt.body)
			req := httptest.NewRequest(tt.method, "/api/addurl", body)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			AddUrl(rr, req)

			responseCode := rr.Code // Получаем статус

			if responseCode != tt.expectedStatus {
				t.Errorf("ожидался статус %d, получен %d", tt.expectedStatus, responseCode)
				return
			}

			if responseCode != 200 {
				return
			}

			var responseJSON map[string]string
			err := json.Unmarshal(rr.Body.Bytes(), &responseJSON)
			if err != nil {
				t.Fatalf("не удалось распарсить JSON: %v", err)
			}

			if !strings.HasPrefix(responseJSON["shorturl"], "http://localhost:3030/") {
				t.Errorf("получен не верный shorturl %s", responseJSON["shorturl"])
			}

		})
	}
}

func TestGetUrl(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		quary          string
		inmemory       bool
		expectedStatus int
	}{
		{"Стандарт", http.MethodGet, `?shorturl=https://example.com`, true, 400},
		{"Неправильный метод", http.MethodPost, `?shorturl=https://example.com`, true, 405},
		{"Неправильный query", http.MethodGet, `?url=https://example.com`, true, 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.inmemory {
				InMemory = true
			}

			req := httptest.NewRequest(tt.method, "/api/heturl"+tt.quary, nil)
			rr := httptest.NewRecorder()
			GetUrl(rr, req)

			responseCode := rr.Code // Получаем статус

			if responseCode != tt.expectedStatus {
				t.Errorf("ожидался статус %d, получен %d", tt.expectedStatus, responseCode)
				return
			}

			if responseCode != 200 {
				return
			}

			var responseJSON map[string]string
			err := json.Unmarshal(rr.Body.Bytes(), &responseJSON)
			if err != nil {
				t.Fatalf("не удалось распарсить JSON: %v", err)
			}

		})
	}
}
