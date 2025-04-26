package repositories

import (
	"context"

	"1337b04rd/internal/domain/models"
)

// UserRepository представляет интерфейс для работы с хранилищем пользователей
type UserRepository interface {
	// GetByID возвращает пользователя по его ID
	GetByID(ctx context.Context, id int64) (*models.User, error)

	// GetBySessionID возвращает пользователя по идентификатору сессии
	GetBySessionID(ctx context.Context, sessionID string) (*models.User, error)

	// Create создает нового пользователя
	Create(ctx context.Context, user *models.User) (int64, error)

	// CreateWithSession создает нового пользователя с сессией
	CreateWithSession(ctx context.Context, user *models.User, sessionID string) (int64, error)

	// GetRandomAvatar получает случайный аватар для пользователя
	GetRandomAvatar(ctx context.Context) (string, error)
}
