package main

import (
	"example.com/m/pkg/api"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("файл .env не найден, используются переменные окружения")
	}

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

	if err := http.ListenAndServe(":3030", apiMux); err != nil {
		log.Fatal("Ошибка сервера:", err)
	}
}
