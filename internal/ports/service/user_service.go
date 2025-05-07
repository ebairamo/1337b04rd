package service

import (
	"context"

	"1337b04rd/internal/domain/models"
)

// UserService определяет интерфейс для сервиса пользователей
type UserService interface {
	// GetByID возвращает пользователя по ID
	GetByID(ctx context.Context, id int64) (*models.User, error)

	// GetUserBySessionID возвращает пользователя по идентификатору сессии
	GetUserBySessionID(ctx context.Context, sessionID string) (*models.User, error)

	// CreateAnonymousUser создает анонимного пользователя
	CreateAnonymousUser(ctx context.Context) (*models.User, error)

	// CreateAnonymousUserWithSession создает анонимного пользователя с сессией
	CreateAnonymousUserWithSession(ctx context.Context, sessionID string) (*models.User, error)
}
