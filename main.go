package main

import (
	"log"
	"net/http"

	"myproject/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/about", handlers.About)
	mux.HandleFunc("/ping", handlers.Ping)
	mux.HandleFunc("/", handlers.Index)
	mux.HandleFunc("/add", handlers.Add)

	// Статика из папки web/static
	mux.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
