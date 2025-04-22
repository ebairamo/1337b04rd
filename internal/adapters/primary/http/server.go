package http

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func StartServer(port int, db *sql.DB) {
	// Создаем новый маршрутизатор
	mux := http.NewServeMux()

	// Регистрируем все маршруты
	RegisterRoutes(mux, db)

	log.Println("Handlers registered")

	// Запуск сервера
	addr := fmt.Sprintf(":%d", port)
	log.Println("Server starting on address:", addr)
	slog.Info("Server started", "address", addr)

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Println("Server failed to start:", err)
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
