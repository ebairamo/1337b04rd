package postgres_test

import (
	"context"
	"database/sql"
	"testing"

	"1337b04rd/internal/adapters/secondary/postgres"
	"1337b04rd/internal/domain/models"
)

// MockAvatarService имитирует сервис аватаров для тестирования
type MockAvatarService struct {
	avatarURLs []string
	nameURLs   []string
	currentIdx int
}

// NewMockAvatarService создает новый экземпляр мок-сервиса аватаров
func NewMockAvatarService() *MockAvatarService {
	return &MockAvatarService{
		avatarURLs: []string{"https://example.com/avatar1.jpg", "https://example.com/avatar2.jpg"},
		nameURLs:   []string{"Test User 1", "Test User 2"},
		currentIdx: 0,
	}
}

// GetRandomAvatar возвращает случайный URL аватара и имя
func (m *MockAvatarService) GetRandomAvatar(ctx context.Context) (string, string, error) {
	if len(m.avatarURLs) == 0 {
		return "https://example.com/default.jpg", "Anonymous", nil
	}
	idx := m.currentIdx % len(m.avatarURLs)
	m.currentIdx++
	return m.avatarURLs[idx], m.nameURLs[idx], nil
}

// ResetUsedIDs сбрасывает список использованных ID
func (m *MockAvatarService) ResetUsedIDs() {
	m.currentIdx = 0
}

// MockDB имитирует БД для тестирования
type MockDB struct {
	users    map[int64]*models.User
	sessions map[string]int64
	lastID   int64
}

// NewMockDB создает новый экземпляр мок-БД
func NewMockDB() *MockDB {
	return &MockDB{
		users:    make(map[int64]*models.User),
		sessions: make(map[string]int64),
		lastID:   0,
	}
}

// QueryRow имитирует sql.QueryRow
func (m *MockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	// Для простоты тестирования считаем, что все запросы на получение пользователя по ID
	// В реальном тесте нужно будет проверять и разбирать SQL запросы
	id, ok := args[0].(int64)
	if !ok {
		return &sql.Row{}
	}
	_, exists := m.users[id]
	if !exists {
		return &sql.Row{}
	}
	// Конвертируем пользователя в результат запроса
	// В реальном тесте это делается через драйвер БД
	return &sql.Row{}
}

// QueryContext имитирует sql.Query
func (m *MockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	// Имитируем результат запроса
	return nil, nil
}

// Exec имитирует sql.Exec
func (m *MockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	// Имитируем результат запроса
	return nil, nil
}

// BeginTx имитирует sql.BeginTx
func (m *MockDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	// Имитируем транзакцию
	return nil, nil
}

// Close имитирует sql.Close
func (m *MockDB) Close() error {
	return nil
}

// Ping имитирует sql.Ping
func (m *MockDB) Ping() error {
	return nil
}

// TestNewUserRepository проверяет создание нового репозитория
func TestNewUserRepository(t *testing.T) {
	// В реальном тесте нужно будет настроить тестовую БД или мок
	// Для этого теста достаточно проверить, что репозиторий создается
	mockDB := &sql.DB{}
	mockAvatarService := NewMockAvatarService()

	repo := postgres.NewUserRepository(mockDB, mockAvatarService)

	if repo == nil {
		t.Fatalf("Не удалось создать репозиторий пользователей")
	}
}

// TestGetRandomAvatar проверяет получение случайного аватара
func TestGetRandomAvatar(t *testing.T) {
	// Создаем мок-сервис аватаров
	mockAvatarService := NewMockAvatarService()

	// Создаем репозиторий с мок-сервисом
	repo := postgres.NewUserRepository(&sql.DB{}, mockAvatarService)

	// Получаем аватар
	avatarURL, err := repo.GetRandomAvatar(context.Background())

	// Проверка результатов
	if err != nil {
		t.Fatalf("Ошибка при получении аватара: %v", err)
	}
	if avatarURL == "" {
		t.Errorf("Получен пустой URL аватара")
	}
	if avatarURL != mockAvatarService.avatarURLs[0] {
		t.Errorf("Неверный URL аватара: ожидалось '%s', получено '%s'", mockAvatarService.avatarURLs[0], avatarURL)
	}
}
