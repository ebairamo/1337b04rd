package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"1337b04rd/internal/adapters/primary/http/middleware"
	"1337b04rd/internal/domain/services"
)

// CommentHandler обрабатывает HTTP запросы для комментариев
type CommentHandler struct {
	commentService *services.CommentService
	userService    *services.UserService
}

// NewCommentHandler создает новый обработчик комментариев
func NewCommentHandler(commentService *services.CommentService, userService *services.UserService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		userService:    userService,
	}
}

// HandleGetComment обрабатывает GET запрос для получения комментария
func (h *CommentHandler) HandleGetComment(w http.ResponseWriter, r *http.Request) {
	// Получаем пользователя из контекста
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		slog.Error("Пользователь не найден в контексте")
		http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
		return
	}

	// Парсим ID комментария из пути
	path := strings.TrimPrefix(r.URL.Path, "/api/comments/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		slog.Error("Невозможно преобразовать ID в число", "path", path, "error", err)
		http.Error(w, "Неверный ID комментария", http.StatusBadRequest)
		return
	}

	comment, err := h.commentService.GetCommentByID(r.Context(), id)
	if err != nil {
		slog.Error("Ошибка получения комментария", "id", id, "error", err)
		http.Error(w, "Комментарий не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

// HandleGetPostComments обрабатывает GET запрос для получения комментариев к посту
func (h *CommentHandler) HandleGetPostComments(w http.ResponseWriter, r *http.Request) {
	// Получаем пользователя из контекста
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		slog.Error("Пользователь не найден в контексте")
		http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
		return
	}

	// Парсим ID поста из пути
	path := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	path = strings.TrimSuffix(path, "/comments")
	postID, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		slog.Error("Невозможно преобразовать ID в число", "path", path, "error", err)
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	// Параметры пагинации
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // По умолчанию
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // По умолчанию
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	comments, err := h.commentService.GetCommentsByPostID(r.Context(), postID, limit, offset)
	if err != nil {
		slog.Error("Ошибка получения комментариев", "post_id", postID, "error", err)
		http.Error(w, "Не удалось получить комментарии", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// HandleCreateComment обрабатывает POST запрос для создания комментария
func (h *CommentHandler) HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		slog.Error("Неверный метод", "method", r.Method)
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Логирование для отладки
	slog.Info("Получен запрос на создание комментария", "path", r.URL.Path, "method", r.Method)

	// Получаем пользователя из контекста
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		slog.Error("Пользователь не найден в контексте")
		http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
		return
	}

	// Парсим данные формы
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 МБ
		slog.Error("Ошибка парсинга multipart формы", "error", err)
		if err = r.ParseForm(); err != nil {
			slog.Error("Ошибка парсинга формы", "error", err)
			http.Error(w, "Ошибка обработки данных формы", http.StatusBadRequest)
			return
		}
	}

	// Получаем ID поста
	postIDStr := r.FormValue("post_id")
	slog.Info("Получен ID поста из формы", "post_id", postIDStr)

	if postIDStr == "" {
		slog.Error("ID поста не предоставлен")
		http.Error(w, "Необходимо указать ID поста", http.StatusBadRequest)
		return
	}

	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		slog.Error("Невозможно преобразовать ID поста в число", "post_id", postIDStr, "error", err)
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	// Получаем текст комментария
	content := r.FormValue("comment")
	slog.Info("Получен текст комментария", "content_length", len(content))

	if content == "" {
		slog.Error("Текст комментария пуст")
		http.Error(w, "Комментарий не может быть пустым", http.StatusBadRequest)
		return
	}

	// Получаем ID родительского комментария (если это ответ)
	replyToIDStr := r.FormValue("reply_to_id")
	var replyToID int64 = 0
	if replyToIDStr != "" {
		replyToID, err = strconv.ParseInt(replyToIDStr, 10, 64)
		if err != nil {
			slog.Error("Невозможно преобразовать ID родительского комментария в число", "reply_to_id", replyToIDStr, "error", err)
			// Не возвращаем ошибку, просто игнорируем reply_to_id
			replyToID = 0
		}
	}

	// Получаем файл изображения (если есть)
	var imageURL string
	file, handler, err := r.FormFile("file")
	if err == nil && file != nil {
		defer file.Close()
		slog.Info("Файл получен", "filename", handler.Filename, "size", handler.Size)

		// TODO: Сохранить файл в S3 и получить URL
		// Здесь должна быть реализация сохранения файла
		imageURL = "https://rickandmortyapi.com/api/character/avatar/1.jpeg" // Заглушка
	}

	// Создаем комментарий через сервис
	comment, err := h.commentService.CreateComment(r.Context(), postID, user.ID, content, imageURL, replyToID)
	if err != nil {
		slog.Error("Ошибка создания комментария", "error", err)
		http.Error(w, "Не удалось создать комментарий: "+err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Комментарий успешно создан", "comment_id", comment.ID, "post_id", postID)

	// Перенаправляем на страницу поста
	http.Redirect(w, r, "/post/"+strconv.FormatInt(postID, 10), http.StatusSeeOther)
}

// HandleDeleteComment обрабатывает DELETE запрос для удаления комментария
func (h *CommentHandler) HandleDeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Получаем пользователя из контекста
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		slog.Error("Пользователь не найден в контексте")
		http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
		return
	}

	// Парсим ID комментария из пути
	path := strings.TrimPrefix(r.URL.Path, "/api/comments/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		slog.Error("Невозможно преобразовать ID в число", "path", path, "error", err)
		http.Error(w, "Неверный ID комментария", http.StatusBadRequest)
		return
	}

	err = h.commentService.DeleteComment(r.Context(), id)
	if err != nil {
		slog.Error("Ошибка удаления комментария", "id", id, "error", err)
		http.Error(w, "Не удалось удалить комментарий", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
