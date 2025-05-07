package rickandmorty

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"
)

const (
	BaseURL        = "https://rickandmortyapi.com/api"
	MaxCharacterID = 826 // По состоянию на 2023 год в API доступно 826 персонажей
)

// AvatarService представляет сервис для работы с API Rick and Morty
type AvatarService struct {
	client   *http.Client
	usedIDs  map[int]bool
	maxRetry int
	baseURL  string // для тестирования
	maxID    int    // для тестирования
}

// Character представляет персонажа из API Rick and Morty
type Character struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	Species string `json:"species"`
	Status  string `json:"status"`
}

// CharacterResponse представляет ответ API для одного персонажа
type CharacterResponse struct {
	Character
	Error string `json:"error"`
}

// NewAvatarService создает новый экземпляр сервиса аватаров
func NewAvatarService() *AvatarService {
	return &AvatarService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		usedIDs:  make(map[int]bool),
		maxRetry: 5,
	}
}

// NewAvatarServiceWithBaseURL создает новый экземпляр сервиса аватаров с указанным URL API
// Используется для тестирования
func NewAvatarServiceWithBaseURL(baseURL string, maxID int) *AvatarService {
	return &AvatarService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		usedIDs:  make(map[int]bool),
		maxRetry: 5,
		baseURL:  baseURL,
		maxID:    maxID,
	}
}

// GetRandomAvatar возвращает URL случайного аватара
func (s *AvatarService) GetRandomAvatar(ctx context.Context) (string, string, error) {
	// Инициализируем генератор случайных чисел
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Используем baseURL и maxID из структуры, если они установлены, иначе используем константы
	baseURL := BaseURL
	if s.baseURL != "" {
		baseURL = s.baseURL
	}

	maxID := MaxCharacterID
	if s.maxID > 0 {
		maxID = s.maxID
	}

	// Пробуем получить аватар несколько раз в случае ошибки
	for i := 0; i < s.maxRetry; i++ {
		// Генерируем случайный ID персонажа
		characterID := rnd.Intn(maxID) + 1

		// Проверяем, не был ли уже использован этот ID
		if len(s.usedIDs) < maxID && s.usedIDs[characterID] {
			continue // Если ID уже использовался, пробуем другой
		}

		// Формируем URL для запроса
		url := fmt.Sprintf("%s/character/%d", baseURL, characterID)

		// Создаем новый HTTP запрос
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			slog.Error("Ошибка создания запроса", "error", err)
			continue
		}

		// Выполняем запрос
		resp, err := s.client.Do(req)
		if err != nil {
			slog.Error("Ошибка выполнения запроса", "error", err)
			continue
		}
		defer resp.Body.Close()

		// Проверяем код ответа
		if resp.StatusCode != http.StatusOK {
			slog.Error("Неудачный ответ API", "status", resp.StatusCode)
			continue
		}

		// Декодируем ответ
		var character CharacterResponse
		if err := json.NewDecoder(resp.Body).Decode(&character); err != nil {
			slog.Error("Ошибка декодирования ответа", "error", err)
			continue
		}

		// Если получили ошибку в ответе, пробуем снова
		if character.Error != "" {
			slog.Error("Ошибка API", "error", character.Error)
			continue
		}

		// Отмечаем ID как использованный
		s.usedIDs[characterID] = true

		// Возвращаем URL аватара и имя персонажа
		return character.Image, character.Name, nil
	}

	// Если все попытки закончились неудачей, возвращаем стандартный аватар
	slog.Warn("Не удалось получить аватар, используем стандартный")
	return "https://rickandmortyapi.com/api/character/avatar/1.jpeg", "Anonymous", nil
}

// ResetUsedIDs сбрасывает список использованных ID
// Это может понадобиться, если количество пользователей превысит количество персонажей
func (s *AvatarService) ResetUsedIDs() {
	s.usedIDs = make(map[int]bool)
}
