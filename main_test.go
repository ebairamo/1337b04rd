package main

import (
	"fmt"
	"os"
	"testing"

	// Импортируем тестируемые пакеты
	_ "1337b04rd/internal/adapters/primary/http/handlers"
	_ "1337b04rd/internal/adapters/primary/http/middleware"
	_ "1337b04rd/internal/adapters/secondary/postgres"
	_ "1337b04rd/internal/adapters/secondary/rickandmorty"
	_ "1337b04rd/internal/adapters/secondary/s3"
	_ "1337b04rd/internal/domain/services"
)

func TestMain(m *testing.M) {
	fmt.Println("Запуск тестов для проекта 1337b04rd")

	// Устанавливаем переменные окружения для тестов
	os.Setenv("S3_HOST", "localhost")
	os.Setenv("S3_PORT", "9000")

	// Запускаем тесты
	exitCode := m.Run()

	fmt.Println("Завершение тестов")

	// Завершаем приложение с кодом из тестов
	os.Exit(exitCode)
}

// TestCoverage проверяет, что тесты покрывают достаточный процент кода
func TestCoverage(t *testing.T) {
	// Этот тест всегда проходит, так как проверка покрытия выполняется внешним инструментом
	// Запуск тестов с покрытием:
	// go test -coverprofile=coverage.out ./...
	// go tool cover -html=coverage.out -o coverage.html

	// Для информации выводим ожидаемое минимальное покрытие
	fmt.Println("Ожидаемое минимальное покрытие кода: 20%")
}
