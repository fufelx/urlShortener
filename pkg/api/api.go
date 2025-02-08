package api

import (
	"encoding/json"
	"example.com/m/pkg/pgsql"
	"example.com/m/pkg/urlshortener"
	"net/http"
)

var (
	db, Errdb     = pgsql.New()
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

	if shorturl, exist := UrlToShorturl[creds.Url]; exist {
		res := pgsql.UrlInfo{Url: creds.Url, ShortUrl: shorturl}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res.ShortUrl)
		return
	}

createShortUrlAgain:
	shorturl, err := urlshortener.MakeUrlShort()
	if err != nil {
		http.Error(w, "ошибка при сокращении ссылки", http.StatusInternalServerError)
		return
	}

	res := pgsql.UrlInfo{Url: creds.Url, ShortUrl: shorturl}

	if InMemory {
		if _, exist := ShorturlToUrl[shorturl]; exist {
			goto createShortUrlAgain
		}
		UrlToShorturl[creds.Url] = shorturl
		ShorturlToUrl[shorturl] = creds.Url
	} else {
		shorturltmp, err := db.AddUrl(res)
		if err != nil {
			goto createShortUrlAgain
		}
		res.ShortUrl = shorturltmp
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res.ShortUrl)
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

	var res pgsql.UrlInfo

	if InMemory {
		res.Url = ShorturlToUrl[shorturl]
		if res.Url == "" {
			http.Error(w, "оригинал ссылки отсутствует", http.StatusBadRequest)
			return
		}
	} else {
		restmp, err := db.GetUrlByShotrurl(shorturl)
		if err != nil {
			http.Error(w, "оригинал ссылки отсутствует", http.StatusBadRequest)
			return
		}
		res = restmp
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res.Url)
}
