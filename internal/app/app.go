package app

import (
	"1337b04rd/internal/adapters/secondary/postgres"
	"log"
	"log/slog"
	"os"
)

func Start(port int) {

	// Настройка логгера
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	log.Println("Logger set up")
	logger.Info("Starting 1337b04rd server", "port", *port)

	db := postgres.Connect()
	defer db.Close()
	logger.Info("База данных подключена")

	
}
