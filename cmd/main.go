package main

import (
	"log"
	"net/http"

	"forum/internal/handler"
	"forum/internal/service"
	"forum/internal/storage"
)

func main() {
	db := storage.InitDB()
	storages := storage.NewStorage(db)
	services := service.NewService(storages)
	handlers := handler.NewHandler(services)
	handlers.InitRoutes()
	log.Println("Запуск веб-сервера на http://localhost:8080")
	http.ListenAndServe(":8080", handlers.Mux)
}
