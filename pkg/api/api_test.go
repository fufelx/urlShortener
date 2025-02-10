package api

import (
	"bytes"
	"encoding/json"
	"example.com/m/pkg/pgsql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockDBfunction struct {
	MockAddUrl           func(info pgsql.UrlInfo) (string, error)
	MockGetUrlByShotrurl func(shorturl string) (pgsql.UrlInfo, error)
}

func (m *MockDBfunction) AddUrl(info pgsql.UrlInfo) (string, error) {
	return m.MockAddUrl(info)
}

func (m *MockDBfunction) GetUrlByShotrurl(shorturl string) (pgsql.UrlInfo, error) {
	return m.MockGetUrlByShotrurl(shorturl)
}

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
		{"Стандарт", http.MethodPost, `{"url": "https://example.com"}`, false, 200},
		{"Неправильный метод", http.MethodGet, `{"url": "https://example.com"}`, false, 405},
		{"Неправильный JSONx1", http.MethodPost, `{"urlshort": "https://example.com"}`, false, 400},
		{"Неправильный URL", http.MethodPost, `{"url": "example.com"}`, false, 400},
		{"Неправильный JSONx2", http.MethodPost, `{"url": "https://example.com"`, false, 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.inmemory {
				InMemory = true
			}

			mockDb := &MockDBfunction{
				MockAddUrl: func(info pgsql.UrlInfo) (string, error) {
					return "http://localhost:3030/hHh33JwTCL", nil
				},
			}

			handler := DB{DBF: DBfunction(mockDb)}

			body := bytes.NewBufferString(tt.body)
			req := httptest.NewRequest(tt.method, "/api/addurl", body)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.AddUrl(rr, req)

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
		{"Стандарт", http.MethodGet, `?shorturl=https://example.com`, false, 200},
		{"Неправильный метод", http.MethodPost, `?shorturl=https://example.com`, false, 405},
		{"Неправильный query", http.MethodGet, `?url=https://example.com`, false, 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.inmemory {
				InMemory = true
			}

			mockDb := &MockDBfunction{
				MockGetUrlByShotrurl: func(shorturl string) (pgsql.UrlInfo, error) {
					return pgsql.UrlInfo{}, nil
				},
			}

			handler := DB{DBF: DBfunction(mockDb)}

			req := httptest.NewRequest(tt.method, "/api/heturl"+tt.quary, nil)
			rr := httptest.NewRecorder()
			handler.GetUrl(rr, req)

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
