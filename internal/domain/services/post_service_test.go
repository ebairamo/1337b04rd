package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/domain/services"
)

// MockPostRepository имитирует репозиторий постов для тестирования
type MockPostRepository struct {
	posts      map[int64]*models.Post
	currentID  int64
	archiveErr error
}

// NewMockPostRepository создает новый экземпляр мок-репозитория
func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{
		posts:     make(map[int64]*models.Post),
		currentID: 1,
	}
}

// GetByID возвращает пост по ID
func (m *MockPostRepository) GetByID(ctx context.Context, id int64) (*models.Post, error) {
	post, exists := m.posts[id]
	if !exists {
		return nil, fmt.Errorf("пост с ID %d не найден", id)
	}
	return post, nil
}

// GetAll возвращает список постов
func (m *MockPostRepository) GetAll(ctx context.Context, limit, offset int, archived bool) ([]*models.Post, error) {
	var result []*models.Post
	for _, post := range m.posts {
		if post.IsArchived == archived {
			result = append(result, post)
		}
	}
	// Упрощенная пагинация для теста
	if offset >= len(result) {
		return []*models.Post{}, nil
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end], nil
}

// GetAllForArchiving возвращает все неархивированные посты
func (m *MockPostRepository) GetAllForArchiving(ctx context.Context) ([]*models.Post, error) {
	var result []*models.Post
	for _, post := range m.posts {
		if !post.IsArchived {
			result = append(result, post)
		}
	}
	return result, nil
}

// Create создает новый пост
func (m *MockPostRepository) Create(ctx context.Context, post *models.Post) (int64, error) {
	id := m.currentID
	post.ID = id
	m.posts[id] = post
	m.currentID++
	return id, nil
}

// Archive архивирует пост
func (m *MockPostRepository) Archive(ctx context.Context, id int64) error {
	if m.archiveErr != nil {
		return m.archiveErr
	}
	post, exists := m.posts[id]
	if !exists {
		return nil
	}
	post.IsArchived = true
	return nil
}

// Тесты для сервиса постов
func TestCreatePost(t *testing.T) {
	// Инициализация мок-репозиториев
	mockPostRepo := NewMockPostRepository()
	mockUserRepo := NewMockUserRepository()

	// Создаем тестового пользователя
	user := &models.User{
		ID:        1,
		Username:  "testuser",
		AvatarURL: "https://example.com/avatar.jpg",
		CreatedAt: time.Now(),
	}
	mockUserRepo.users[user.ID] = user

	// Инициализация сервиса
	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	// Создаем пост
	title := "Test Post"
	content := "This is a test post content"
	imageURL := "https://example.com/image.jpg"

	post, err := postService.CreatePost(context.Background(), title, content, imageURL, user.ID)

	// Проверка результатов
	if err != nil {
		t.Fatalf("Ошибка при создании поста: %v", err)
	}
	if post == nil {
		t.Fatalf("Пост не создан")
	}
	if post.ID != 1 {
		t.Errorf("Неверный ID поста: ожидалось 1, получено %d", post.ID)
	}
	if post.Title != title {
		t.Errorf("Неверный заголовок поста: ожидалось '%s', получено '%s'", title, post.Title)
	}
	if post.Content != content {
		t.Errorf("Неверное содержимое поста: ожидалось '%s', получено '%s'", content, post.Content)
	}
	if post.ImageURL != imageURL {
		t.Errorf("Неверный URL изображения: ожидалось '%s', получено '%s'", imageURL, post.ImageURL)
	}
	if post.UserID != user.ID {
		t.Errorf("Неверный ID пользователя: ожидалось %d, получено %d", user.ID, post.UserID)
	}
	if post.UserName != user.Username {
		t.Errorf("Неверное имя пользователя: ожидалось '%s', получено '%s'", user.Username, post.UserName)
	}
	if post.IsArchived {
		t.Errorf("Пост создан как архивный, хотя должен быть неархивным")
	}
}

func TestGetPostByID(t *testing.T) {
	// Инициализация мок-репозиториев
	mockPostRepo := NewMockPostRepository()
	mockUserRepo := NewMockUserRepository()

	// Создаем тестовый пост
	testPost := &models.Post{
		ID:        1,
		Title:     "Test Post",
		Content:   "This is a test post content",
		ImageURL:  "https://example.com/image.jpg",
		UserID:    1,
		UserName:  "testuser",
		AvatarURL: "https://example.com/avatar.jpg",
		CreatedAt: time.Now(),
	}
	mockPostRepo.posts[testPost.ID] = testPost

	// Инициализация сервиса
	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	// Получаем пост по ID
	post, err := postService.GetPostByID(context.Background(), testPost.ID)

	// Проверка результатов
	if err != nil {
		t.Fatalf("Ошибка при получении поста по ID: %v", err)
	}
	if post == nil {
		t.Fatalf("Пост не найден")
	}
	if post.ID != testPost.ID {
		t.Errorf("Неверный ID поста: ожидалось %d, получено %d", testPost.ID, post.ID)
	}
	if post.Title != testPost.Title {
		t.Errorf("Неверный заголовок поста: ожидалось '%s', получено '%s'", testPost.Title, post.Title)
	}

	// Тест на получение несуществующего поста
	post, err = postService.GetPostByID(context.Background(), 999)
	if err == nil {
		t.Errorf("Ожидалась ошибка при получении несуществующего поста")
	}
	if post != nil {
		t.Errorf("Пост получен, хотя не должен был")
	}
}

func TestArchivePost(t *testing.T) {
	// Инициализация мок-репозиториев
	mockPostRepo := NewMockPostRepository()
	mockUserRepo := NewMockUserRepository()

	// Создаем тестовый пост
	testPost := &models.Post{
		ID:         1,
		Title:      "Test Post",
		Content:    "This is a test post content",
		UserID:     1,
		UserName:   "testuser",
		CreatedAt:  time.Now(),
		IsArchived: false,
	}
	mockPostRepo.posts[testPost.ID] = testPost

	// Инициализация сервиса
	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	// Архивируем пост
	err := postService.ArchivePost(context.Background(), testPost.ID)

	// Проверка результатов
	if err != nil {
		t.Fatalf("Ошибка при архивации поста: %v", err)
	}
	if !testPost.IsArchived {
		t.Errorf("Пост не был архивирован")
	}
}
