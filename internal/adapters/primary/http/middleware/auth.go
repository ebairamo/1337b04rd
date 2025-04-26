package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/domain/services"
)

const (
	cookieName   = "session_id"
	cookieMaxAge = 7 * 24 * 60 * 60 // 1 неделя в секундах
)

// Ключ контекста для хранения пользователя
type userContextKey string

const UserContextKey userContextKey = "user"

// AuthMiddleware представляет собой middleware для аутентификации
type AuthMiddleware struct {
	userService *services.UserService
}

// NewAuthMiddleware создает новый экземпляр middleware аутентификации
func NewAuthMiddleware(userService *services.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
	}
}

// Handler обрабатывает аутентификацию пользователя
func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем или создаем сессию
		user, err := m.getOrCreateSession(r, w)
		if err != nil {
			slog.Error("Ошибка при работе с сессией", "error", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}

		// Добавляем пользователя в контекст запроса
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getOrCreateSession получает существующего пользователя или создает нового
func (m *AuthMiddleware) getOrCreateSession(r *http.Request, w http.ResponseWriter) (*models.User, error) {
	// Ищем куки
	cookie, err := r.Cookie(cookieName)

	// Если куки нет или произошла ошибка, создаем нового пользователя
	if err != nil || cookie.Value == "" {
		return m.createSession(w)
	}

	// Пытаемся получить пользователя по идентификатору сессии
	sessionID := cookie.Value
	user, err := m.userService.GetUserBySessionID(r.Context(), sessionID)

	// Если пользователь не найден или произошла ошибка, создаем нового
	if err != nil || user == nil {
		return m.createSession(w)
	}

	return user, nil
}

// createSession создает нового пользователя и устанавливает куки
func (m *AuthMiddleware) createSession(w http.ResponseWriter) (*models.User, error) {
	// Генерируем UUID для сессии
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	// Создаем анонимного пользователя
	user, err := m.userService.CreateAnonymousUserWithSession(context.Background(), sessionID)
	if err != nil {
		return nil, err
	}

	// Устанавливаем куки
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   cookieMaxAge,
		Expires:  time.Now().Add(time.Duration(cookieMaxAge) * time.Second),
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)

	return user, nil
}

// GetUserFromContext извлекает пользователя из контекста
func GetUserFromContext(ctx context.Context) *models.User {
	user, ok := ctx.Value(UserContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}

// generateSessionID генерирует случайный идентификатор сессии в формате UUID
func generateSessionID() (string, error) {
	// Генерируем 16 случайных байт для UUID
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	// Устанавливаем биты версии (4) и варианта (2) согласно RFC4122
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // версия 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // вариант 2

	// Форматируем UUID в строку стандартного формата
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]), nil
}
