package api

import (
	"encoding/json"
	"example.com/m/pkg/pgsql"
	"example.com/m/pkg/urlshortener"
	"net/http"
	"net/url"
)

var (
	InMemory      = false
	UrlToShorturl = make(map[string]string) // мапа для хранения оригинальной ссылки и её сокращения
	ShorturlToUrl = make(map[string]string) // мапа для хранения сокращения и её оригинальной ссылки
)

func AddUrl(w http.ResponseWriter, r *http.Request) {
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

	type answer struct {
		Url string `json:"shorturl"`
	}

	// Если ссылка уже существует, возвращаем её
	if shorturl, exist := UrlToShorturl[creds.Url]; exist {
		ans := answer{Url: shorturl}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ans)
		return
	}

createShortUrlAgain:
	shorturl := urlshortener.MakeUrlShort()
	res := pgsql.UrlInfo{Url: creds.Url, ShortUrl: shorturl}
	ans := answer{Url: shorturl}
	if InMemory {
		if _, exist := ShorturlToUrl[shorturl]; exist {
			// Если сгенерированная короткая ссылка уже существует, генерируем новую
			goto createShortUrlAgain
		}
		UrlToShorturl[creds.Url] = shorturl
		ShorturlToUrl[shorturl] = creds.Url
	} else {
		shorturltmp, err := pgsql.Db.AddUrl(res)
		if err != nil {
			// Если сгенерированная короткая ссылка уже существует, генерируем новую
			goto createShortUrlAgain
		}
		ans.Url = shorturltmp
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ans)
}

func GetUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "неправильный метод", http.StatusMethodNotAllowed)
		return
	}

	shorturl := r.URL.Query().Get("shorturl")
	if shorturl == "" {
		http.Error(w, "shorturl не найдена в RawQuery", http.StatusBadRequest)
		return
	}

	type answer struct {
		Url string `json:"url"`
	}
	ans := answer{}
	if InMemory {
		ans.Url = ShorturlToUrl[shorturl]
		if ans.Url == "" {
			http.Error(w, "оригинал ссылки отсутствует", http.StatusBadRequest)
			return
		}
	} else {
		restmp, err := pgsql.Db.GetUrlByShotrurl(shorturl)
		if err != nil {
			http.Error(w, "оригинал ссылки отсутствует", http.StatusBadRequest)
			return
		}
		ans.Url = restmp.Url
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ans)
}
