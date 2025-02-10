package redirect

import (
	"example.com/m/pkg/api"
	"net/http"
)

func RedirectUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "неправильный метод", http.StatusMethodNotAllowed)
		return
	}
	shorturl := "http://localhost:3030/" + r.URL.Path[1:]

	if shorturl == "" {
		http.Error(w, "shorturl не найден", http.StatusBadRequest)
		return
	}

	var originalUrl string

	if api.InMemory {
		// Проверяем в памяти
		originalUrl = api.ShorturlToUrl[shorturl]
	} else {
		// Проверяем в базе данных
		res, err := api.Db.GetUrlByShotrurl(shorturl)
		if err != nil {
			http.Error(w, "shorturl не существует", http.StatusNotFound)
			return
		}
		originalUrl = res.Url
	}

	// Если оригинальный URL найден, делаем редирект
	if originalUrl != "" {
		http.Redirect(w, r, originalUrl, http.StatusFound) // 302 Redirect
		return
	}

	http.Error(w, "shorturl не найден", http.StatusNotFound)
}
