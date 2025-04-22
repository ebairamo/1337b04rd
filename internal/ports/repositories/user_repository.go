package repositories

import (
	"1337b04rd/internal/domain/models"
	"context"
)

// UserRepository представляет интерфейс для работы с хранилищем пользователей
type UserRepository interface {
	// GetByID возвращает пользователя по его ID
	GetByID(ctx context.Context, id int64) (*models.User, error)

	// Create создает нового пользователя
	Create(ctx context.Context, user *models.User) (int64, error)

	// GetRandomAvatar получает случайный аватар для пользователя
	GetRandomAvatar(ctx context.Context) (string, error)
}
