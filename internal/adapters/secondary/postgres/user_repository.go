package postgres

import (
	"1337b04rd/internal/domain/models"
	"context"
	"database/sql"
	"log/slog"
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
	slog.Info("Заглушка: получение пользователя по ID", "id", id)
	return &models.User{
		ID:        id,
		Username:  "anonymous",
		AvatarURL: "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
	}, nil
}

// Create создает нового пользователя
func (r *UserRepository) Create(ctx context.Context, user *models.User) (int64, error) {
	// TODO: реализовать создание пользователя в БД
	slog.Info("Заглушка: создание пользователя", "username", user.Username)
	return 1, nil
}

// GetRandomAvatar получает случайный аватар для пользователя
func (r *UserRepository) GetRandomAvatar(ctx context.Context) (string, error) {
	// TODO: реализовать получение случайного аватара
	return "https://rickandmortyapi.com/api/character/avatar/1.jpeg", nil
}
