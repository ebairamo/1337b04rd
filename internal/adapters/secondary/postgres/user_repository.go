package postgres

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/ports/external"
)

// UserRepository реализует интерфейс репозитория пользователей для PostgreSQL
type UserRepository struct {
	db            *sql.DB
	avatarService external.AvatarService
}

// NewUserRepository создает новый экземпляр репозитория пользователей
func NewUserRepository(db *sql.DB, avatarService external.AvatarService) *UserRepository {
	return &UserRepository{
		db:            db,
		avatarService: avatarService,
	}
}

// GetByID возвращает пользователя по его ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `SELECT id, user_name, avatar_url, created_at 
			  FROM users 
			  WHERE id = $1`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.AvatarURL, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Пользователь не найден", "id", id)
			return nil, err
		}
		slog.Error("Ошибка при получении пользователя из БД", "error", err)
		return nil, err
	}

	slog.Info("Пользователь найден", "user", user)
	return &user, nil
}

// GetBySessionID возвращает пользователя по идентификатору сессии
func (r *UserRepository) GetBySessionID(ctx context.Context, sessionID string) (*models.User, error) {
	query := `SELECT u.id, u.user_name, u.avatar_url, u.created_at 
			  FROM users u
			  JOIN sessions s ON u.id = s.user_id
			  WHERE s.id = $1 AND s.expires_at > NOW()`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&user.ID, &user.Username, &user.AvatarURL, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Сессия не найдена или истекла", "session_id", sessionID)
			return nil, err
		}
		slog.Error("Ошибка при получении пользователя по сессии", "error", err)
		return nil, err
	}

	slog.Info("Пользователь найден по сессии", "user_id", user.ID)
	return &user, nil
}

// Create создает нового пользователя
func (r *UserRepository) Create(ctx context.Context, user *models.User) (int64, error) {
	query := `INSERT INTO users (user_name, avatar_url, created_at)
			  VALUES ($1, $2, $3)
			  RETURNING id`

	var id int64
	err := r.db.QueryRowContext(ctx, query,
		user.Username, user.AvatarURL, time.Now()).Scan(&id)
	if err != nil {
		slog.Error("Ошибка при создании пользователя", "error", err)
		return 0, err
	}

	slog.Info("Пользователь создан", "id", id)
	return id, nil
}

// CreateWithSession создает нового пользователя с сессией
func (r *UserRepository) CreateWithSession(ctx context.Context, user *models.User, sessionID string) (int64, error) {
	// Начинаем транзакцию
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("Ошибка начала транзакции", "error", err)
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Создаем пользователя
	userQuery := `INSERT INTO users (user_name, avatar_url, created_at)
				 VALUES ($1, $2, $3)
				 RETURNING id`

	var userID int64
	err = tx.QueryRowContext(ctx, userQuery,
		user.Username, user.AvatarURL, time.Now()).Scan(&userID)
	if err != nil {
		slog.Error("Ошибка при создании пользователя в транзакции", "error", err)
		return 0, err
	}

	// Создаем сессию
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 1 неделя
	sessionQuery := `INSERT INTO sessions (id, user_id, avatar_url, created_at, expires_at)
					VALUES ($1, $2, $3, $4, $5)`

	_, err = tx.ExecContext(ctx, sessionQuery,
		sessionID, userID, user.AvatarURL, time.Now(), expiresAt)
	if err != nil {
		slog.Error("Ошибка при создании сессии", "error", err)
		return 0, err
	}

	// Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		slog.Error("Ошибка при коммите транзакции", "error", err)
		return 0, err
	}

	slog.Info("Пользователь и сессия созданы", "user_id", userID, "session_id", sessionID)
	return userID, nil
}

// GetRandomAvatar получает случайный аватар для пользователя
func (r *UserRepository) GetRandomAvatar(ctx context.Context) (string, error) {
	// Используем сервис аватаров Rick and Morty
	avatarURL, _, err := r.avatarService.GetRandomAvatar(ctx)
	if err != nil {
		slog.Error("Ошибка получения аватара из Rick and Morty API", "error", err)
		return "https://rickandmortyapi.com/api/character/avatar/1.jpeg", nil
	}

	return avatarURL, nil
}
