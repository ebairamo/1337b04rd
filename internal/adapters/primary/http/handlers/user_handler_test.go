package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"1337b04rd/internal/adapters/primary/http/handlers"
	"1337b04rd/internal/adapters/primary/http/middleware"
	"1337b04rd/internal/domain/models"
)

// MockUserService имитирует сервис пользователей для тестирования
type MockUserService struct {
	users        map[int64]*models.User
	sessionUsers map[string]*models.User
	currentID    int64
}

// NewMockUserService создает новый экземпляр мок-сервиса
func NewMockUserService() *MockUserService {
	return &MockUserService{
		users:        make(map[int64]*models.User),
		sessionUsers: make(map[string]*models.User),
		currentID:    1,
	}
}

// GetByID возвращает пользователя по ID
func (m *MockUserService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, nil
	}
	return user, nil
}

// GetUserBySessionID возвращает пользователя по ID сессии
func (m *MockUserService) GetUserBySessionID(ctx context.Context, sessionID string) (*models.User, error) {
	user, exists := m.sessionUsers[sessionID]
	if !exists {
		return nil, nil
	}
	return user, nil
}

// CreateAnonymousUser создает анонимного пользователя
func (m *MockUserService) CreateAnonymousUser(ctx context.Context) (*models.User, error) {
	user := &models.User{
		ID:        m.currentID,
		Username:  "anonymous",
		AvatarURL: "https://example.com/avatar.jpg",
		CreatedAt: time.Now(),
	}
	m.users[user.ID] = user
	m.currentID++
	return user, nil
}

// CreateAnonymousUserWithSession создает анонимного пользователя с сессией
func (m *MockUserService) CreateAnonymousUserWithSession(ctx context.Context, sessionID string) (*models.User, error) {
	user, err := m.CreateAnonymousUser(ctx)
	if err != nil {
		return nil, err
	}
	m.sessionUsers[sessionID] = user
	return user, nil
}

// TestHandleGetUser тестирует обработчик для получения пользователя
func TestHandleGetUser(t *testing.T) {
	// Инициализация мок-сервиса
	mockService := NewMockUserService()

	// Создаем тестового пользователя
	testUser := &models.User{
		ID:        1,
		Username:  "testuser",
		AvatarURL: "https://example.com/avatar.jpg",
		CreatedAt: time.Now(),
	}
	mockService.users[testUser.ID] = testUser

	// Создаем обработчик
	handler := handlers.NewUserHandler(mockService)

	// Создаем тестовый HTTP запрос
	req, err := http.NewRequest("GET", "/api/users/1", nil)
	if err != nil {
		t.Fatalf("Ошибка создания запроса: %v", err)
	}

	// Устанавливаем пользователя в контекст запроса
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, testUser)
	req = req.WithContext(ctx)

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	handler.HandleGetUser(rr, req)

	// Проверяем статус ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неверный статус: ожидалось %v, получено %v", http.StatusOK, status)
	}

	// Проверяем заголовок Content-Type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Неверный Content-Type: ожидалось application/json, получено %v", contentType)
	}

	// Декодируем ответ
	var responseUser models.User
	if err := json.Unmarshal(rr.Body.Bytes(), &responseUser); err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}

	// Проверяем данные пользователя
	if responseUser.ID != testUser.ID {
		t.Errorf("Неверный ID: ожидалось %v, получено %v", testUser.ID, responseUser.ID)
	}
	if responseUser.Username != testUser.Username {
		t.Errorf("Неверное имя: ожидалось %v, получено %v", testUser.Username, responseUser.Username)
	}
	if responseUser.AvatarURL != testUser.AvatarURL {
		t.Errorf("Неверный URL аватара: ожидалось %v, получено %v", testUser.AvatarURL, responseUser.AvatarURL)
	}
}

// TestHandleCreateUser тестирует обработчик для создания пользователя
func TestHandleCreateUser(t *testing.T) {
	// Инициализация мок-сервиса
	mockService := NewMockUserService()

	// Создаем обработчик
	handler := handlers.NewUserHandler(mockService)

	// Создаем тестовый HTTP запрос
	req, err := http.NewRequest("POST", "/api/users", nil)
	if err != nil {
		t.Fatalf("Ошибка создания запроса: %v", err)
	}

	// Создаем тестового пользователя для контекста
	contextUser := &models.User{
		ID:        999,
		Username:  "contextuser",
		AvatarURL: "https://example.com/context-avatar.jpg",
		CreatedAt: time.Now(),
	}

	// Устанавливаем пользователя в контекст запроса
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, contextUser)
	req = req.WithContext(ctx)

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	handler.HandleCreateUser(rr, req)

	// Проверяем статус ответа
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Неверный статус: ожидалось %v, получено %v", http.StatusCreated, status)
	}

	// Проверяем заголовок Content-Type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Неверный Content-Type: ожидалось application/json, получено %v", contentType)
	}

	// Декодируем ответ
	var responseUser models.User
	if err := json.Unmarshal(rr.Body.Bytes(), &responseUser); err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}

	// Проверяем данные пользователя
	if responseUser.ID != 1 {
		t.Errorf("Неверный ID: ожидалось 1, получено %v", responseUser.ID)
	}
	if responseUser.Username != "anonymous" {
		t.Errorf("Неверное имя: ожидалось anonymous, получено %v", responseUser.Username)
	}
	if responseUser.AvatarURL == "" {
		t.Errorf("Пустой URL аватара")
	}
}
