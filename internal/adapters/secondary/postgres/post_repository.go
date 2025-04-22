package postgres

import (
	"1337b04rd/internal/domain/models"
	"context"
	"database/sql"
	"log/slog"
	"time"
)

// PostRepository реализует интерфейс репозитория постов для PostgreSQL
type PostRepository struct {
	db *sql.DB
}

// NewPostRepository создает новый экземпляр репозитория постов
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

// GetByID возвращает пост по его ID
func (r *PostRepository) GetByID(ctx context.Context, id int64) (*models.Post, error) {
	// TODO: реализовать получение поста из БД
	slog.Info("Заглушка: получение поста по ID", "id", id)

	return &models.Post{
		ID:         id,
		Title:      "Тестовый пост",
		Content:    "Это тестовое содержимое поста для заглушки репозитория",
		ImageURL:   "https://www.google.com/images/branding/googlelogo/2x/googlelogo_light_color_272x92dp.png",
		UserID:     1,
		UserName:   "anonymous",
		AvatarURL:  "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
		CreatedAt:  time.Now().Add(-24 * time.Hour), // Создан день назад
		IsArchived: false,
	}, nil
}

// GetAll возвращает все посты с возможной фильтрацией
func (r *PostRepository) GetAll(ctx context.Context, limit, offset int, archived bool) ([]*models.Post, error) {
	// TODO: реализовать получение всех постов из БД
	query := `SELECT * FROM posts WHERE is_archived = $3 ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset, archived)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var posts []*models.Post

	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.ImageURL, &post.UserID, &post.UserName, &post.AvatarURL, &post.CreatedAt, &post.IsArchived)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	if err := rows.Err(); err != nil {
		slog.Error("Ошибка при обработке строк из БД", "error", err)
		return nil, err
	}
	return posts, nil
}

// Create создает новый пост
func (r *PostRepository) Create(ctx context.Context, post *models.Post) (int64, error) {
	// TODO: реализовать создание поста в БД
	slog.Info("Заглушка: создание поста", "title", post.Title)
	return 1, nil
}

// Update обновляет существующий пост
func (r *PostRepository) Update(ctx context.Context, post *models.Post) error {
	// TODO: реализовать обновление поста в БД
	slog.Info("Заглушка: обновление поста", "id", post.ID)
	return nil
}

// Archive архивирует пост
func (r *PostRepository) Archive(ctx context.Context, id int64) error {
	// TODO: реализовать архивацию поста в БД
	slog.Info("Заглушка: архивация поста", "id", id)
	return nil
}
