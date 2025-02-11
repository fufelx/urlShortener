package handlers

import (
	"encoding/json"
	"example.com/m/internal/service"
	"example.com/m/internal/storage"
	"net/http"
	"net/url"
)

type ShortUrlAnswer struct {
	Url string `json:"shorturl"`
}

type OriginalUrlAnswer struct {
	Url string `json:"url"`
}

func AddUrl(db *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "неправильный метод", http.StatusMethodNotAllowed)
			return
		}

		var creds struct {
			Url string `json:"url"`
		}

		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil || creds.Url == "" {
			http.Error(w, "неправильный JSON", http.StatusBadRequest)
			return
		}

		_, err = url.ParseRequestURI(creds.Url)
		if err != nil {
			http.Error(w, "неправильный формат ссылки", http.StatusBadRequest)
			return
		}

		//Формируем ответ
		action, err := service.NewLinik(db, creds.Url)
		if err != nil {
			http.Error(w, "ошибка при обработке ссылки", http.StatusInternalServerError)
			return
		}

		ans := ShortUrlAnswer{Url: "http://localhost:3030/" + action}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ans)
	}
}

func GetUrl(db *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "неправильный метод", http.StatusMethodNotAllowed)
			return
		}

		shorturl := r.URL.Query().Get("shorturl")
		if shorturl == "" {
			http.Error(w, "shorturl не найдена в RawQuery", http.StatusBadRequest)
			return
		}

		action, err := service.GetOriginalLink(db, shorturl)
		if err != nil {
			http.Error(w, "ошибка при обработке ссылки", http.StatusInternalServerError)
			return
		}

		ans := OriginalUrlAnswer{Url: action}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ans)
	}
}

func RedirectUrl(db *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "неправильный метод", http.StatusMethodNotAllowed)
			return
		}
		shorturl := r.URL.Path[1:]

		if shorturl == "" {
			http.Error(w, "shorturl не найден", http.StatusBadRequest)
			return
		}

		action, err := service.GetOriginalLink(db, shorturl)
		if err != nil {
			http.Error(w, "оригинальный URL не найден", http.StatusInternalServerError)
			return
		}

		// Если оригинальный URL найден, делаем редирект
		http.Redirect(w, r, action, http.StatusFound) // 302 Redirect
		return
	}
}
