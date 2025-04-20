package http

import (
	"1337b04rd/internal/adapters/primary/http/handlers"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func StartServer(port int) {
	// Регистрация обработчиков
	http.HandleFunc("/", handlers.HandlePage)

	log.Println("Handler registered")

	// Запуск сервера
	addr := fmt.Sprintf(":%d", port)
	log.Println("Server starting on address:", addr)
	slog.Info("Server started", "address", addr)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Println("Server failed to start:", err)
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
