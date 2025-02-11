package main

import (
	"example.com/m/internal/handlers"
	"example.com/m/internal/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	// Проверяем тип памяти(in-memory или pgsql)
	db := &storage.Store{}
	inMemory := os.Getenv("STORAGE_TYPE")
	if inMemory == "in-memory" {
		storage.InMemory = true
	} else {
		initdb, err := storage.New()
		if err != nil {
			log.Fatal(err)
		}
		db = initdb
	}

	apiMux := http.NewServeMux()
	// Добавляем маршруты для API
	apiMux.HandleFunc("/api/addurl", handlers.AddUrl(db))
	apiMux.HandleFunc("/api/geturl", handlers.GetUrl(db))
	apiMux.HandleFunc("/", handlers.RedirectUrl(db))

	if err := http.ListenAndServe(":3030", apiMux); err != nil {
		log.Fatal("Ошибка сервера:", err)
	}
}
