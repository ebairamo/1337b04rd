package repositories

import (
	"1337b04rd/internal/domain/models"
	"context"
)

// CommentRepository представляет интерфейс для работы с хранилищем комментариев
type CommentRepository interface {
	// GetByID возвращает комментарий по его ID
	GetByID(ctx context.Context, id int64) (*models.Comment, error)

	// GetByPostID возвращает все комментарии к указанному посту
	GetByPostID(ctx context.Context, postID int64, limit, offset int) ([]*models.Comment, error)

	// Create создает новый комментарий
	Create(ctx context.Context, comment *models.Comment) (int64, error)

	// Delete удаляет комментарий по ID
	Delete(ctx context.Context, id int64) error
}
