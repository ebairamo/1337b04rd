package services

import (
	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/ports/repositories"
	"context"
	"log/slog"
	"time"
)

// UserService предоставляет бизнес-логику для работы с пользователями
type UserService struct {
	userRepo repositories.UserRepository
}

// NewUserService создает новый экземпляр сервиса пользователей
func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetByID возвращает пользователя по ID
func (s *UserService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	slog.Info("Получение пользователя", "id", id)
	return s.userRepo.GetByID(ctx, id)
}

// CreateAnonymousUser создает анонимного пользователя
func (s *UserService) CreateAnonymousUser(ctx context.Context) (*models.User, error) {
	avatarURL, err := s.userRepo.GetRandomAvatar(ctx)
	if err != nil {
		slog.Error("Ошибка получения аватара", "error", err)
		avatarURL = "https://rickandmortyapi.com/api/character/avatar/1.jpeg"
	}

	user := &models.User{
		Username:  "anonymous",
		AvatarURL: avatarURL,
		CreatedAt: time.Now(),
	}

	id, err := s.userRepo.Create(ctx, user)
	if err != nil {
		slog.Error("Ошибка создания пользователя", "error", err)
		return nil, err
	}

	user.ID = id
	return user, nil
}