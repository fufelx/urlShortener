package main

import (
	"example.com/m/pkg/api"
	"example.com/m/pkg/redirect"
	"log"
	"net/http"
	"os"
)

func main() {
	// Проверяем тип памяти(in-memory или pgsql)
	inMemory := os.Getenv("STORAGE_TYPE")
	if inMemory == "in-memory" {
		api.InMemory = true
	} else {
		if api.Errdb != nil {
			log.Fatal(api.Errdb)
		}
	}

	apiMux := http.NewServeMux()
	// Добавляем маршруты для API
	apiMux.HandleFunc("/api/addurl", api.AddUrl)
	apiMux.HandleFunc("/api/geturl", api.GetUrl)
	apiMux.HandleFunc("/", redirect.RedirectUrl)

	if err := http.ListenAndServe(":3030", apiMux); err != nil {
		log.Fatal("Ошибка сервера:", err)
	}
}
