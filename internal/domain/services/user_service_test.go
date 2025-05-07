package services_test

import (
	"context"
	"testing"
	"time"

	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/domain/services"
)

// MockUserRepository имитирует репозиторий пользователей для тестирования
type MockUserRepository struct {
	users       map[int64]*models.User
	sessions    map[string]int64
	avatarURLs  []string
	currentID   int64
	createFunc  func(ctx context.Context, user *models.User) (int64, error)
	sessionFunc func(ctx context.Context, user *models.User, sessionID string) (int64, error)
}

// NewMockUserRepository создает новый экземпляр мок-репозитория
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:      make(map[int64]*models.User),
		sessions:   make(map[string]int64),
		avatarURLs: []string{"https://example.com/avatar1.jpg", "https://example.com/avatar2.jpg"},
		currentID:  1,
	}
}

// GetByID реализация метода интерфейса для получения пользователя по ID
func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, nil
	}
	return user, nil
}

// GetBySessionID реализация метода для получения пользователя по SessionID
func (m *MockUserRepository) GetBySessionID(ctx context.Context, sessionID string) (*models.User, error) {
	userID, exists := m.sessions[sessionID]
	if !exists {
		return nil, nil
	}
	return m.GetByID(ctx, userID)
}

// Create реализация метода для создания пользователя
func (m *MockUserRepository) Create(ctx context.Context, user *models.User) (int64, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	id := m.currentID
	user.ID = id
	m.users[id] = user
	m.currentID++
	return id, nil
}

// CreateWithSession реализация метода для создания пользователя с сессией
func (m *MockUserRepository) CreateWithSession(ctx context.Context, user *models.User, sessionID string) (int64, error) {
	if m.sessionFunc != nil {
		return m.sessionFunc(ctx, user, sessionID)
	}
	id, err := m.Create(ctx, user)
	if err != nil {
		return 0, err
	}
	m.sessions[sessionID] = id
	return id, nil
}

// GetRandomAvatar реализация метода для получения случайного аватара
func (m *MockUserRepository) GetRandomAvatar(ctx context.Context) (string, error) {
	if len(m.avatarURLs) == 0 {
		return "https://example.com/default.jpg", nil
	}
	return m.avatarURLs[0], nil
}

func TestCreateAnonymousUser(t *testing.T) {
	// Инициализация мок-репозитория
	mockRepo := NewMockUserRepository()
	// Инициализация сервиса с мок-репозиторием
	userService := services.NewUserService(mockRepo)

	// Создание анонимного пользователя через сервис
	user, err := userService.CreateAnonymousUser(context.Background())

	// Проверка результатов
	if err != nil {
		t.Fatalf("Ошибка при создании анонимного пользователя: %v", err)
	}
	if user == nil {
		t.Fatalf("Пользователь не создан")
	}
	if user.ID != 1 {
		t.Errorf("Неверный ID пользователя: ожидалось 1, получено %d", user.ID)
	}
	if user.Username != "anonymous" {
		t.Errorf("Неверное имя пользователя: ожидалось 'anonymous', получено '%s'", user.Username)
	}
	if user.AvatarURL == "" {
		t.Errorf("URL аватара не установлен")
	}
}

func TestGetUserBySessionID(t *testing.T) {
	// Инициализация мок-репозитория
	mockRepo := NewMockUserRepository()

	// Создаем тестового пользователя и сессию
	testUser := &models.User{
		Username:  "testuser",
		AvatarURL: "https://example.com/avatar.jpg",
		CreatedAt: time.Now(),
	}

	// Добавляем пользователя и сессию в мок-репозиторий
	mockRepo.CreateWithSession(context.Background(), testUser, "test-session-id")

	// Инициализация сервиса с мок-репозиторием
	userService := services.NewUserService(mockRepo)

	// Тест на получение пользователя по существующему sessionID
	user, err := userService.GetUserBySessionID(context.Background(), "test-session-id")
	if err != nil {
		t.Fatalf("Ошибка при получении пользователя по sessionID: %v", err)
	}
	if user == nil {
		t.Fatalf("Пользователь не найден по sessionID")
	}
	if user.Username != "testuser" {
		t.Errorf("Неверное имя пользователя: ожидалось 'testuser', получено '%s'", user.Username)
	}

	// Тест на получение пользователя по несуществующему sessionID
	user, err = userService.GetUserBySessionID(context.Background(), "non-existent-id")
	if err != nil {
		t.Logf("Ожидаемая ошибка при несуществующем sessionID: %v", err)
	}
	if user != nil {
		t.Errorf("Пользователь найден по несуществующему sessionID")
	}
}
