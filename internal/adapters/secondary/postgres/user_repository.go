package postgres

import (
	"1337b04rd/internal/domain/models"
	"context"
	"database/sql"
	"log/slog"
	"time"
)

// UserRepository реализует интерфейс репозитория пользователей для PostgreSQL
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository создает новый экземпляр репозитория пользователей
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetByID возвращает пользователя по его ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	// TODO: реализовать получение пользователя из БД

	query := `SELECT 
	user_id, 
	user_name, 
	avatar_url,
	created_at
	FROM posts 
	WHERE user_id = $1`

	var user models.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.AvatarURL, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Пользователь не найден", "id", id)
		}
		slog.Error("Ошибка при получении пользователя из БД", "error", err)
	}
	slog.Info("Пользователь найден", "user", user)
	return &user, nil
}

// Create создает нового пользователя
func (r *UserRepository) Create(ctx context.Context, user *models.User) (int64, error) {
	currentTime := time.Now()

	query := `INSERT INTO users (user_name, avatar_url, created_at)
	VALUES ($1, $2, $3)
	RETURNING user_id`

	err := r.db.QueryRowContext(ctx, query, user.Username, user.AvatarURL, currentTime).Scan(&user.ID)
	if err != nil {
		slog.Error("Ошибка при создании пользователя в БД", "error", err)
		return 0, err
	}
	slog.Info("Пользователь создан", "user", user)
	return user.ID, nil
}

// GetRandomAvatar получает случайный аватар для пользователя
func (r *UserRepository) GetRandomAvatar(ctx context.Context) (string, error) {
	// TODO: реализовать получение случайного аватара
	return "https://rickandmortyapi.com/api/character/avatar/1.jpeg", nil
}
