package rickandmorty_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"1337b04rd/internal/adapters/secondary/rickandmorty"
)

// TestGetRandomAvatar проверяет получение случайного аватара
func TestGetRandomAvatar(t *testing.T) {
	// Создаем мок-сервер для имитации API Rick and Morty
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что запрос направлен на API персонажей
		if !strings.HasPrefix(r.URL.Path, "/api/character/") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Отправляем мок-ответ с данными персонажа
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Получаем ID персонажа из пути
		idStr := strings.TrimPrefix(r.URL.Path, "/api/character/")

		// Создаем разные тестовые данные для разных ID
		var responseData map[string]interface{}
		switch idStr {
		case "1":
			responseData = map[string]interface{}{
				"id":      1,
				"name":    "Rick Sanchez",
				"image":   "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
				"species": "Human",
				"status":  "Alive",
			}
		case "2":
			responseData = map[string]interface{}{
				"id":      2,
				"name":    "Morty Smith",
				"image":   "https://rickandmortyapi.com/api/character/avatar/2.jpeg",
				"species": "Human",
				"status":  "Alive",
			}
		default:
			// Для любого другого ID
			responseData = map[string]interface{}{
				"id":      idStr,
				"name":    "Test Character",
				"image":   "https://rickandmortyapi.com/api/character/avatar/" + idStr + ".jpeg",
				"species": "Test Species",
				"status":  "Test Status",
			}
		}

		// Отправляем ответ
		json.NewEncoder(w).Encode(responseData)
	}))
	defer server.Close()

	// Создаем сервис аватаров с указанием базового URL
	avatarService := rickandmorty.NewAvatarServiceWithBaseURL(server.URL, 826)

	// Тестируем получение случайного аватара
	avatarURL, name, err := avatarService.GetRandomAvatar(context.Background())

	// Проверка результатов
	if err != nil {
		t.Fatalf("Ошибка при получении аватара: %v", err)
	}
	if avatarURL == "" {
		t.Errorf("Получен пустой URL аватара")
	}
	if name == "" {
		t.Errorf("Получено пустое имя персонажа")
	}

	// Проверяем, что URL аватара имеет ожидаемый формат
	if !strings.HasPrefix(avatarURL, "https://rickandmortyapi.com/api/character/avatar/") {
		t.Errorf("Неверный формат URL аватара: %s", avatarURL)
	}
}

// TestResetUsedIDs проверяет сброс использованных ID
func TestResetUsedIDs(t *testing.T) {
	// Создаем мок-сервер для имитации API Rick and Morty
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Отправляем мок-ответ с данными персонажа
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		responseData := map[string]interface{}{
			"id":      1,
			"name":    "Test Character",
			"image":   "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
			"species": "Test Species",
			"status":  "Test Status",
		}

		json.NewEncoder(w).Encode(responseData)
	}))
	defer server.Close()

	// Создаем сервис аватаров с ограниченным количеством персонажей (10) и тестовым URL
	maxCharacters := 10
	avatarService := rickandmorty.NewAvatarServiceWithBaseURL(server.URL, maxCharacters)

	// Получаем максимальное количество аватаров
	for i := 0; i < maxCharacters; i++ {
		avatarService.GetRandomAvatar(context.Background())
	}

	// Сбрасываем использованные ID
	avatarService.ResetUsedIDs()

	// Проверяем, что теперь можно снова получить аватары
	avatarURL, name, err := avatarService.GetRandomAvatar(context.Background())
	if err != nil {
		t.Fatalf("Ошибка при получении аватара после сброса: %v", err)
	}
	if avatarURL == "" {
		t.Errorf("Получен пустой URL аватара после сброса")
	}
	if name == "" {
		t.Errorf("Получено пустое имя персонажа после сброса")
	}
}
