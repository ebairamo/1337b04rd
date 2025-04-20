package postgres

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	connStr := "user=elite_user password=elite_pass dbname=eliteboard_db host=db port=5432 sslmode=disable"
	slog.Info("Строка подключения", "connStr", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("Ошибка подключения к БД", "error", err)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		slog.Error("Не удалось проверить подключение", "error", err)
		os.Exit(1)
	}

	slog.Info("Успешное подключение к базе данных")
	return db
}
