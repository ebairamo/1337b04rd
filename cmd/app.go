package cmd

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	httpAdapter "1337b04rd/internal/adapters/primary/http"
	"1337b04rd/internal/adapters/secondary/postgres"
)

func Run() {
	// Используем обычный log пакет для гарантированного вывода
	log.Println("Starting program...")

	// Парсинг командной строки
	port := flag.Int("port", 8080, "Port number")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	log.Println("Parsed flags, port:", *port, "help:", *help)

	// Показать помощь и выйти, если запрошена помощь
	if *help {
		fmt.Println("hacker board")
		fmt.Println("\nUsage:")
		fmt.Println("  1337b04rd [--port <N>]")
		fmt.Println("  1337b04rd --help")
		fmt.Println("\nOptions:")
		fmt.Println("  --help       Show this screen.")
		fmt.Println("  --port N     Port number.")
		os.Exit(0)
	}

	// Настройка логгера
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// Подключение к базе данных
	db := postgres.Connect()
	defer db.Close()
	logger.Info("База данных подключена")

	// Запускаем HTTP сервер с передачей порта и соединения с БД
	httpAdapter.StartServer(*port, db)
}
