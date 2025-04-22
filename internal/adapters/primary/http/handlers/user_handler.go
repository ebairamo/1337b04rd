package handlers

import (
	"1337b04rd/internal/domain/services"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

// UserHandler обрабатывает HTTP запросы для пользователей
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler создает новый обработчик пользователей
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// HandleGetUser обрабатывает GET запрос для получения пользователя
func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	// Простой парсинг пути для извлечения ID пользователя
	// Пример: /api/users/123
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		slog.Error("Невозможно преобразовать ID в число", "path", path, "error", err)
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		slog.Error("Ошибка получения пользователя", "id", id, "error", err)
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// HandleCreateUser обрабатывает POST запрос для создания пользователя
func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	// В этой заглушке мы просто создаем анонимного пользователя
	user, err := h.userService.CreateAnonymousUser(r.Context())
	if err != nil {
		slog.Error("Ошибка создания пользователя", "error", err)
		http.Error(w, "Не удалось создать пользователя", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
