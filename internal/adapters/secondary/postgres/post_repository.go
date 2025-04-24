package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"1337b04rd/internal/domain/models"
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
	query := `SELECT 
        id, title, content, image_url, user_id, user_name, avatar_url, created_at, is_archived 
        FROM posts 
        WHERE id = $1`

	var post models.Post
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID, &post.Title, &post.Content, &post.ImageURL,
		&post.UserID, &post.UserName, &post.AvatarURL,
		&post.CreatedAt, &post.IsArchived)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Пост не найден", "id", id)
			return nil, fmt.Errorf("post with id %d not found", id)
		}
		slog.Error("Ошибка получения поста", "id", id, "error", err)
		return nil, err
	}

	// Добавим более подробное логирование
	slog.Info("Пост получен",
		"id", post.ID,
		"title", post.Title,
		"content_length", len(post.Content),
		"image_url", post.ImageURL,
		"user_id", post.UserID,
		"user_name", post.UserName,
		"created_at", post.CreatedAt,
		"is_archived", post.IsArchived)

	return &post, nil
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
	currentTime := time.Now()

	query := `INSERT INTO posts (title, content, image_url, user_id, user_name, avatar_url, created_at, is_archived)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id
	`
	var newID int64
	err := r.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ImageURL, post.UserID, post.UserName, post.AvatarURL, currentTime, post.IsArchived).Scan(&newID)
	if err != nil {
		slog.Error("Ошибка создания поста", "error", err)
		return 0, err
	}
	slog.Info("Пост создан",
		"id", newID,
		"title", post.Title,
		"content_length", len(post.Content),
		"image_url", post.ImageURL,
		"user_id", post.UserID,
		"user_name", post.UserName,
		"created_at", currentTime,
		"is_archived", post.IsArchived)
	return newID, err
}

// Archive архивирует пост
func (r *PostRepository) Archive(ctx context.Context, id int64) error {
	// TODO: реализовать архивацию поста в БД
	slog.Info("Заглушка: архивация поста", "id", id)
	return nil
}
