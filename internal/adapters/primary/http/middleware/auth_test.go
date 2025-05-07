package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"1337b04rd/internal/adapters/primary/http/middleware"
	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/ports/service"
)

// MockUserService имитирует сервис пользователей для тестирования middleware аутентификации
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

// Убедимся, что MockUserService реализует интерфейс service.UserService
var _ service.UserService = (*MockUserService)(nil)

// TestAuthMiddlewareNewSession тестирует создание новой сессии
func TestAuthMiddlewareNewSession(t *testing.T) {
	// Инициализация мок-сервиса
	mockService := NewMockUserService()

	// Создаем middleware
	authMiddleware := middleware.NewAuthMiddleware(mockService)

	// Создаем простой обработчик для тестирования
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что пользователь добавлен в контекст
		user := middleware.GetUserFromContext(r.Context())
		if user == nil {
			t.Errorf("Пользователь не добавлен в контекст")
			http.Error(w, "Пользователь не найден", http.StatusInternalServerError)
			return
		}

		// Проверяем данные пользователя
		if user.ID != 1 {
			t.Errorf("Неверный ID пользователя: ожидалось 1, получено %d", user.ID)
		}
		if user.Username != "anonymous" {
			t.Errorf("Неверное имя пользователя: ожидалось 'anonymous', получено '%s'", user.Username)
		}

		w.WriteHeader(http.StatusOK)
	})

	// Оборачиваем тестовый обработчик в middleware
	handler := authMiddleware.Handler(testHandler)

	// Создаем тестовый HTTP запрос без cookie
	req := httptest.NewRequest("GET", "/test", nil)

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	handler.ServeHTTP(rr, req)

	// Проверяем статус ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неверный статус: ожидалось %v, получено %v", http.StatusOK, status)
	}

	// Проверяем, что cookie установлен
	cookies := rr.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatalf("Cookie не установлен")
	}

	sessionCookie := cookies[0]
	if sessionCookie.Name != "session_id" {
		t.Errorf("Неверное имя cookie: ожидалось 'session_id', получено '%s'", sessionCookie.Name)
	}
	if sessionCookie.Value == "" {
		t.Errorf("Пустое значение session_id")
	}
	if !sessionCookie.HttpOnly {
		t.Errorf("Cookie должен быть HttpOnly")
	}
}

// TestAuthMiddlewareExistingSession тестирует работу с существующей сессией
func TestAuthMiddlewareExistingSession(t *testing.T) {
	// Инициализация мок-сервиса
	mockService := NewMockUserService()

	// Создаем тестового пользователя и сессию
	testUser := &models.User{
		ID:        42,
		Username:  "testuser",
		AvatarURL: "https://example.com/avatar42.jpg",
		CreatedAt: time.Now(),
	}
	testSessionID := "test-session-id"
	mockService.users[testUser.ID] = testUser
	mockService.sessionUsers[testSessionID] = testUser

	// Создаем middleware
	authMiddleware := middleware.NewAuthMiddleware(mockService)

	// Создаем простой обработчик для тестирования
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что пользователь добавлен в контекст
		user := middleware.GetUserFromContext(r.Context())
		if user == nil {
			t.Errorf("Пользователь не добавлен в контекст")
			http.Error(w, "Пользователь не найден", http.StatusInternalServerError)
			return
		}

		// Проверяем данные пользователя
		if user.ID != testUser.ID {
			t.Errorf("Неверный ID пользователя: ожидалось %d, получено %d", testUser.ID, user.ID)
		}
		if user.Username != testUser.Username {
			t.Errorf("Неверное имя пользователя: ожидалось '%s', получено '%s'", testUser.Username, user.Username)
		}

		w.WriteHeader(http.StatusOK)
	})

	// Оборачиваем тестовый обработчик в middleware
	handler := authMiddleware.Handler(testHandler)

	// Создаем тестовый HTTP запрос с cookie
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: testSessionID,
	})

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	handler.ServeHTTP(rr, req)

	// Проверяем статус ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неверный статус: ожидалось %v, получено %v", http.StatusOK, status)
	}
}

// TestGetUserFromContext тестирует получение пользователя из контекста
func TestGetUserFromContext(t *testing.T) {
	// Создаем тестового пользователя
	testUser := &models.User{
		ID:        123,
		Username:  "testuser",
		AvatarURL: "https://example.com/avatar.jpg",
		CreatedAt: time.Now(),
	}

	// Создаем контекст с пользователем
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, testUser)

	// Получаем пользователя из контекста
	user := middleware.GetUserFromContext(ctx)

	// Проверяем, что пользователь получен корректно
	if user == nil {
		t.Fatalf("Пользователь не получен из контекста")
	}
	if user.ID != testUser.ID {
		t.Errorf("Неверный ID пользователя: ожидалось %d, получено %d", testUser.ID, user.ID)
	}
	if user.Username != testUser.Username {
		t.Errorf("Неверное имя пользователя: ожидалось '%s', получено '%s'", testUser.Username, user.Username)
	}

	// Проверяем получение пользователя из пустого контекста
	emptyCtx := context.Background()
	emptyUser := middleware.GetUserFromContext(emptyCtx)
	if emptyUser != nil {
		t.Errorf("Получен пользователь из пустого контекста")
	}
}
