package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/domain/services"
)

// MockCommentRepository имитирует репозиторий комментариев для тестирования
type MockCommentRepository struct {
	comments      map[int64]*models.Comment
	postComments  map[int64][]*models.Comment
	currentID     int64
	lastCommentID int64
}

// NewMockCommentRepository создает новый экземпляр мок-репозитория комментариев
func NewMockCommentRepository() *MockCommentRepository {
	return &MockCommentRepository{
		comments:     make(map[int64]*models.Comment),
		postComments: make(map[int64][]*models.Comment),
		currentID:    1,
	}
}

// GetByID возвращает комментарий по ID
func (m *MockCommentRepository) GetByID(ctx context.Context, id int64) (*models.Comment, error) {
	comment, exists := m.comments[id]
	if !exists {
		return nil, fmt.Errorf("комментарий с ID %d не найден", id)
	}
	return comment, nil
}

// GetByPostID возвращает комментарии к посту
func (m *MockCommentRepository) GetByPostID(ctx context.Context, postID int64, limit, offset int) ([]*models.Comment, error) {
	comments, exists := m.postComments[postID]
	if !exists {
		return []*models.Comment{}, nil
	}

	// Упрощенная пагинация для теста
	if offset >= len(comments) {
		return []*models.Comment{}, nil
	}
	end := offset + limit
	if end > len(comments) {
		end = len(comments)
	}
	return comments[offset:end], nil
}

// GetLastCommentByPostID возвращает последний комментарий к посту
func (m *MockCommentRepository) GetLastCommentByPostID(ctx context.Context, postID int64) (*models.Comment, error) {
	comments, exists := m.postComments[postID]
	if !exists || len(comments) == 0 {
		return nil, nil
	}
	return comments[len(comments)-1], nil
}

// Create создает новый комментарий
func (m *MockCommentRepository) Create(ctx context.Context, comment *models.Comment) (int64, error) {
	id := m.currentID
	comment.ID = id
	m.comments[id] = comment

	// Добавляем комментарий в список комментариев к посту
	if _, exists := m.postComments[comment.PostID]; !exists {
		m.postComments[comment.PostID] = []*models.Comment{}
	}
	m.postComments[comment.PostID] = append(m.postComments[comment.PostID], comment)

	m.currentID++
	m.lastCommentID = id
	return id, nil
}

// Delete удаляет комментарий
func (m *MockCommentRepository) Delete(ctx context.Context, id int64) error {
	comment, exists := m.comments[id]
	if !exists {
		return fmt.Errorf("комментарий с ID %d не найден", id)
	}

	// Удаляем комментарий из списка комментариев к посту
	if comments, exists := m.postComments[comment.PostID]; exists {
		for i, c := range comments {
			if c.ID == id {
				m.postComments[comment.PostID] = append(comments[:i], comments[i+1:]...)
				break
			}
		}
	}

	delete(m.comments, id)
	return nil
}

func TestCreateComment(t *testing.T) {
	// Инициализация мок-репозиториев
	mockCommentRepo := NewMockCommentRepository()
	mockUserRepo := NewMockUserRepository()
	mockPostRepo := NewMockPostRepository()

	// Создаем тестового пользователя и пост
	user := &models.User{
		ID:        1,
		Username:  "testuser",
		AvatarURL: "https://example.com/avatar.jpg",
		CreatedAt: time.Now(),
	}
	mockUserRepo.users[user.ID] = user

	post := &models.Post{
		ID:        1,
		Title:     "Test Post",
		Content:   "This is a test post content",
		UserID:    user.ID,
		UserName:  user.Username,
		AvatarURL: user.AvatarURL,
		CreatedAt: time.Now(),
	}
	mockPostRepo.posts[post.ID] = post

	// Инициализация сервиса
	commentService := services.NewCommentService(mockCommentRepo, mockUserRepo, mockPostRepo)

	// Создаем комментарий
	content := "This is a test comment"
	imageURL := "https://example.com/comment-image.jpg"

	comment, err := commentService.CreateComment(context.Background(), post.ID, user.ID, content, imageURL, 0)

	// Проверка результатов
	if err != nil {
		t.Fatalf("Ошибка при создании комментария: %v", err)
	}
	if comment == nil {
		t.Fatalf("Комментарий не создан")
	}
	if comment.ID != 1 {
		t.Errorf("Неверный ID комментария: ожидалось 1, получено %d", comment.ID)
	}
	if comment.PostID != post.ID {
		t.Errorf("Неверный ID поста: ожидалось %d, получено %d", post.ID, comment.PostID)
	}
	if comment.UserID != user.ID {
		t.Errorf("Неверный ID пользователя: ожидалось %d, получено %d", user.ID, comment.UserID)
	}
	if comment.Content != content {
		t.Errorf("Неверное содержимое комментария: ожидалось '%s', получено '%s'", content, comment.Content)
	}
	if comment.ImageURL != imageURL {
		t.Errorf("Неверный URL изображения: ожидалось '%s', получено '%s'", imageURL, comment.ImageURL)
	}
	if comment.ReplyToID != 0 {
		t.Errorf("Неверный ID родительского комментария: ожидалось 0, получено %d", comment.ReplyToID)
	}
}

func TestCreateReplyComment(t *testing.T) {
	// Инициализация мок-репозиториев
	mockCommentRepo := NewMockCommentRepository()
	mockUserRepo := NewMockUserRepository()
	mockPostRepo := NewMockPostRepository()

	// Создаем тестовых пользователей и пост
	user1 := &models.User{
		ID:        1,
		Username:  "user1",
		AvatarURL: "https://example.com/avatar1.jpg",
		CreatedAt: time.Now(),
	}
	mockUserRepo.users[user1.ID] = user1

	user2 := &models.User{
		ID:        2,
		Username:  "user2",
		AvatarURL: "https://example.com/avatar2.jpg",
		CreatedAt: time.Now(),
	}
	mockUserRepo.users[user2.ID] = user2

	post := &models.Post{
		ID:        1,
		Title:     "Test Post",
		Content:   "This is a test post content",
		UserID:    user1.ID,
		UserName:  user1.Username,
		AvatarURL: user1.AvatarURL,
		CreatedAt: time.Now(),
	}
	mockPostRepo.posts[post.ID] = post

	// Инициализация сервиса
	commentService := services.NewCommentService(mockCommentRepo, mockUserRepo, mockPostRepo)

	// Создаем первый комментарий
	comment1, err := commentService.CreateComment(context.Background(), post.ID, user1.ID, "First comment", "", 0)
	if err != nil {
		t.Fatalf("Ошибка при создании первого комментария: %v", err)
	}

	// Создаем ответ на первый комментарий
	replyContent := "Reply to first comment"
	reply, err := commentService.CreateComment(context.Background(), post.ID, user2.ID, replyContent, "", comment1.ID)

	// Проверка результатов
	if err != nil {
		t.Fatalf("Ошибка при создании ответа на комментарий: %v", err)
	}
	if reply == nil {
		t.Fatalf("Ответ на комментарий не создан")
	}
	if reply.ReplyToID != comment1.ID {
		t.Errorf("Неверный ID родительского комментария: ожидалось %d, получено %d", comment1.ID, reply.ReplyToID)
	}
	if reply.PostID != post.ID {
		t.Errorf("Неверный ID поста: ожидалось %d, получено %d", post.ID, reply.PostID)
	}
	if reply.UserID != user2.ID {
		t.Errorf("Неверный ID пользователя: ожидалось %d, получено %d", user2.ID, reply.UserID)
	}
	if reply.Content != replyContent {
		t.Errorf("Неверное содержимое ответа: ожидалось '%s', получено '%s'", replyContent, reply.Content)
	}
}

func TestGetCommentsByPostID(t *testing.T) {
	// Инициализация мок-репозиториев
	mockCommentRepo := NewMockCommentRepository()
	mockUserRepo := NewMockUserRepository()
	mockPostRepo := NewMockPostRepository()

	// Создаем тестового пользователя и пост
	user := &models.User{
		ID:        1,
		Username:  "testuser",
		AvatarURL: "https://example.com/avatar.jpg",
		CreatedAt: time.Now(),
	}
	mockUserRepo.users[user.ID] = user

	post := &models.Post{
		ID:        1,
		Title:     "Test Post",
		Content:   "This is a test post content",
		UserID:    user.ID,
		UserName:  user.Username,
		AvatarURL: user.AvatarURL,
		CreatedAt: time.Now(),
	}
	mockPostRepo.posts[post.ID] = post

	// Создаем несколько комментариев к посту
	comment1 := &models.Comment{
		PostID:    post.ID,
		UserID:    user.ID,
		UserName:  user.Username,
		AvatarURL: user.AvatarURL,
		Content:   "Comment 1",
		CreatedAt: time.Now(),
	}
	comment2 := &models.Comment{
		PostID:    post.ID,
		UserID:    user.ID,
		UserName:  user.Username,
		AvatarURL: user.AvatarURL,
		Content:   "Comment 2",
		CreatedAt: time.Now(),
	}

	mockCommentRepo.Create(context.Background(), comment1)
	mockCommentRepo.Create(context.Background(), comment2)

	// Инициализация сервиса
	commentService := services.NewCommentService(mockCommentRepo, mockUserRepo, mockPostRepo)

	// Получаем комментарии к посту
	comments, err := commentService.GetCommentsByPostID(context.Background(), post.ID, 10, 0)

	// Проверка результатов
	if err != nil {
		t.Fatalf("Ошибка при получении комментариев к посту: %v", err)
	}
	if len(comments) != 2 {
		t.Errorf("Неверное количество комментариев: ожидалось 2, получено %d", len(comments))
	}
	if comments[0].Content != "Comment 1" {
		t.Errorf("Неверное содержимое первого комментария: ожидалось 'Comment 1', получено '%s'", comments[0].Content)
	}
	if comments[1].Content != "Comment 2" {
		t.Errorf("Неверное содержимое второго комментария: ожидалось 'Comment 2', получено '%s'", comments[1].Content)
	}

	// Проверка пагинации
	comments, err = commentService.GetCommentsByPostID(context.Background(), post.ID, 1, 1)
	if err != nil {
		t.Fatalf("Ошибка при получении комментариев с пагинацией: %v", err)
	}
	if len(comments) != 1 {
		t.Errorf("Неверное количество комментариев: ожидалось 1, получено %d", len(comments))
	}
	if comments[0].Content != "Comment 2" {
		t.Errorf("Неверное содержимое комментария: ожидалось 'Comment 2', получено '%s'", comments[0].Content)
	}
}
