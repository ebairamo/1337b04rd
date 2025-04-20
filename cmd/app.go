package cmd

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"

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

	// Настройка обработчиков
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request:", r.URL.Path)
		var page string
		name := r.FormValue("name")
		subject := r.FormValue("subject")
		comment := r.FormValue("comment")
		file, fileHeader, err := r.FormFile("file")
		fmt.Println("file:", file, "fileHeader:", fileHeader, "err:", err, "name:", name, "subject:", subject, "comment:", comment)
		Pat := r.URL.Path
		if Pat == "/" {
			page = "templates" + "/" + "catalog.html"
		} else {
			page = "templates" + "/" + Pat
		}

		tmpl, _ := template.ParseFiles(page)
		tmpl.Execute(w, nil)
	})

	log.Println("Handler registered")

	// Запуск сервера
	addr := fmt.Sprintf(":%d", *port)
	log.Println("Server starting on address:", addr)
	logger.Info("Server started", "address", addr)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Println("Server failed to start:", err)
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
