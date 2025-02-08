package api_test

import (
	"bytes"
	"encoding/json"
	"example.com/m/pkg/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddUrl(t *testing.T) {
	api.InMemory = false
	api.UrlToShorturl = make(map[string]string)
	api.ShorturlToUrl = make(map[string]string)

	requestBody, _ := json.Marshal(map[string]string{
		"url": "https://example.com",
	})

	req, err := http.NewRequest("POST", "/add", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(api.AddUrl)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	}

	var shortUrl string
	if err := json.Unmarshal(recorder.Body.Bytes(), &shortUrl); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}

	if len(shortUrl) == 0 {
		t.Error("short URL is empty")
	}
}

func TestGetUrl(t *testing.T) {
	api.InMemory = false
	shortUrl := "http://localhost:3030/hHh33JwTCL"
	originalUrl := "https://example.com"
	api.ShorturlToUrl[shortUrl] = originalUrl

	req, err := http.NewRequest("GET", "/get?shorturl="+shortUrl, nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(api.GetUrl)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	}

	var returnedUrl string
	if err := json.Unmarshal(recorder.Body.Bytes(), &returnedUrl); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}

	if returnedUrl != originalUrl {
		t.Errorf("expected %s, got %s", originalUrl, returnedUrl)
	}
}
