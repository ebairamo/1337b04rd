package postgres

import (
	"1337b04rd/internal/domain/models"
	"context"
	"database/sql"
	"log/slog"
	"time"
)

// CommentRepository реализует интерфейс репозитория комментариев для PostgreSQL
type CommentRepository struct {
	db *sql.DB
}

// NewCommentRepository создает новый экземпляр репозитория комментариев
func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

// GetByID возвращает комментарий по его ID
func (r *CommentRepository) GetByID(ctx context.Context, id int64) (*models.Comment, error) {
	slog.Info("Заглушка: получение комментария", "id", id)
	return &models.Comment{
		ID:        id,
		PostID:    1,
		UserID:    1,
		UserName:  "Anonymous",
		AvatarURL: "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
		Content:   "Тестовый комментарий",
		CreatedAt: time.Now(),
	}, nil
}

// GetByPostID возвращает все комментарии к указанному посту
func (r *CommentRepository) GetByPostID(ctx context.Context, postID int64, limit, offset int) ([]*models.Comment, error) {
	slog.Info("Заглушка: получение комментариев к посту", "post_id", postID)
	// Создаем тестовые комментарии для заглушки
	comments := []*models.Comment{
		{
			ID:        1,
			PostID:    postID,
			UserID:    1,
			UserName:  "Anonymous1",
			AvatarURL: "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
			Content:   "Тестовый комментарий 1",
			CreatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        2,
			PostID:    postID,
			UserID:    2,
			UserName:  "Anonymous2",
			AvatarURL: "https://rickandmortyapi.com/api/character/avatar/2.jpeg",
			Content:   "Тестовый комментарий 2",
			CreatedAt: time.Now(),
		},
	}
	return comments, nil
}

// Create создает новый комментарий
func (r *CommentRepository) Create(ctx context.Context, comment *models.Comment) (int64, error) {
	slog.Info("Заглушка: создание комментария",
		"post_id", comment.PostID,
		"user_id", comment.UserID,
		"content", comment.Content)
	return 1, nil
}

// Delete удаляет комментарий по ID
func (r *CommentRepository) Delete(ctx context.Context, id int64) error {
	slog.Info("Заглушка: удаление комментария", "id", id)
	return nil
}
