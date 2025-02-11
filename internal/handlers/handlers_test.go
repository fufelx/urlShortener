package handlers

import (
	"bytes"
	"encoding/json"
	"example.com/m/internal/storage"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestAddUrl(t *testing.T) {

	var db *storage.Store

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
				storage.InMemory = true
			}

			body := bytes.NewBufferString(tt.body)
			req := httptest.NewRequest(tt.method, "/api/addurl", body)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler := AddUrl(db)
			handler.ServeHTTP(rr, req)

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

	var db *storage.Store

	tests := []struct {
		name           string
		method         string
		quary          string
		inmemory       bool
		expectedStatus int
	}{
		{"Стандарт", http.MethodGet, `?shorturl=example`, true, 500},
		{"Неправильный метод", http.MethodPost, `?shorturl=https://example.com`, true, 405},
		{"Неправильный query", http.MethodGet, `?url=https://example.com`, true, 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.inmemory {
				storage.InMemory = true
			}

			req := httptest.NewRequest(tt.method, "/api/geturl"+tt.quary, nil)
			rr := httptest.NewRecorder()
			handler := GetUrl(db)
			handler.ServeHTTP(rr, req)

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

func TestRedirectUrl(t *testing.T) {
	type args struct {
		db *storage.Store
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RedirectUrl(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedirectUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
