package repositories

import (
	"1337b04rd/internal/domain/models"
	"context"
)

// PostRepository представляет интерфейс для работы с хранилищем постов
type PostRepository interface {
	// GetByID возвращает пост по его ID
	GetByID(ctx context.Context, id int64) (*models.Post, error)

	// GetAll возвращает все посты с возможной фильтрацией
	GetAll(ctx context.Context, limit, offset int, archived bool) ([]*models.Post, error)

	// Create создает новый пост
	Create(ctx context.Context, post *models.Post) (int64, error)
	// Archive архивирует пост
	Archive(ctx context.Context, id int64) error
}
