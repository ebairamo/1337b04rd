package services

import (
	"context"
	"log/slog"
	"time"

	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/ports/repositories"
)

// PostService предоставляет бизнес-логику для работы с постами
type PostService struct {
	postRepo repositories.PostRepository
	userRepo repositories.UserRepository
}

// NewPostService создает новый экземпляр сервиса постов
func NewPostService(postRepo repositories.PostRepository, userRepo repositories.UserRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

// GetPostByID возвращает пост по ID
func (s *PostService) GetPostByID(ctx context.Context, id int64) (*models.Post, error) {
	slog.Info("Получение поста", "id", id)
	return s.postRepo.GetByID(ctx, id)
}

// GetAllPosts возвращает список постов
func (s *PostService) GetAllPosts(ctx context.Context, limit, offset int, archived bool) ([]*models.Post, error) {
	slog.Info("Получение списка постов", "limit", limit, "offset", offset, "archived", archived)

	if limit <= 0 {
		limit = 10 // По умолчанию 10 постов
	}

	return s.postRepo.GetAll(ctx, limit, offset, archived)
}

// GetTotalPostsCount возвращает общее количество постов
func (s *PostService) GetTotalPostsCount(ctx context.Context, archived bool) (int, error) {
	slog.Info("Получение общего количества постов", "archived", archived)
	// В реальной реализации здесь бы был дополнительный метод в репозитории
	// для подсчета общего количества постов

	// Сейчас просто возвращаем фиктивные данные для демонстрации
	if archived {
		return 50, nil // Предположим, что в архиве 50 постов
	}
	return 100, nil // Предположим, что в основном каталоге 100 постов
}

// CreatePost создает новый пост
func (s *PostService) CreatePost(ctx context.Context, title, content, imageURL string, userID int64) (*models.Post, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		slog.Error("Ошибка получения пользователя", "error", err)
		return nil, err
	}

	post := &models.Post{
		Title:      title,
		Content:    content,
		ImageURL:   imageURL,
		UserID:     userID,
		UserName:   user.Username,
		AvatarURL:  user.AvatarURL,
		CreatedAt:  time.Now(),
		IsArchived: false,
	}

	id, err := s.postRepo.Create(ctx, post)
	if err != nil {
		slog.Error("Ошибка создания поста", "error", err)
		return nil, err
	}

	post.ID = id
	return post, nil
}

// ArchivePost архивирует пост
func (s *PostService) ArchivePost(ctx context.Context, id int64) error {
	slog.Info("Архивация поста", "id", id)
	return s.postRepo.Archive(ctx, id)
}
