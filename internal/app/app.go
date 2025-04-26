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
	log.Println("Запуск программы...")

	// Парсинг командной строки
	port := flag.Int("port", 8080, "Номер порта")
	help := flag.Bool("help", false, "Показать справку")
	flag.Parse()

	log.Println("Параметры обработаны, порт:", *port, "справка:", *help)

	// Показать помощь и выйти, если запрошена помощь
	if *help {
		fmt.Println("hacker board")
		fmt.Println("\nИспользование:")
		fmt.Println("  1337b04rd [--port <N>]")
		fmt.Println("  1337b04rd --help")
		fmt.Println("\nОпции:")
		fmt.Println("  --help       Показать эту справку.")
		fmt.Println("  --port N     Номер порта.")
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
