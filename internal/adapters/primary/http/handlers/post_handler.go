package handlers

import (
	"1337b04rd/internal/domain/services"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

// PostHandler обрабатывает HTTP запросы для постов
type PostHandler struct {
	postService *services.PostService
	userService *services.UserService
}

// NewPostHandler создает новый обработчик постов
func NewPostHandler(postService *services.PostService, userService *services.UserService) *PostHandler {
	return &PostHandler{
		postService: postService,
		userService: userService,
	}
}

// HandleGetPost обрабатывает GET запрос для получения поста
func (h *PostHandler) HandleGetPost(w http.ResponseWriter, r *http.Request) {
	// Простой парсинг пути для извлечения ID поста
	// Пример: /api/posts/123
	path := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		slog.Error("Невозможно преобразовать ID в число", "path", path, "error", err)
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	post, err := h.postService.GetPostByID(r.Context(), id)
	if err != nil {
		slog.Error("Ошибка получения поста", "id", id, "error", err)
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	}

	// Проверяем, нужно ли вернуть JSON или HTML
	contentType := r.Header.Get("Accept")
	if strings.Contains(contentType, "application/json") {
		// Возвращаем JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	} else {
		// Возвращаем HTML страницу поста
		tmpl, err := template.ParseFiles("templates/post.html")
		if err != nil {
			slog.Error("Ошибка загрузки шаблона", "error", err)
			http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
			return
		}

		// TODO: Загрузить комментарии
		tmpl.Execute(w, post)
	}
}

// HandleGetAllPosts обрабатывает GET запрос для получения списка постов
func (h *PostHandler) HandleGetAllPosts(w http.ResponseWriter, r *http.Request) {
	// Параметры запроса
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	archivedStr := r.URL.Query().Get("archived")

	limit := 10
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	archived := false
	if archivedStr == "true" || archivedStr == "1" {
		archived = true
	}

	posts, err := h.postService.GetAllPosts(r.Context(), limit, offset, archived)
	if err != nil {
		slog.Error("Ошибка получения списка постов", "error", err)
		http.Error(w, "Не удалось получить список постов", http.StatusInternalServerError)
		return
	}

	// Проверяем, нужно ли вернуть JSON или HTML
	contentType := r.Header.Get("Accept")
	if strings.Contains(contentType, "application/json") {
		// Возвращаем JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	} else {
		// Возвращаем HTML страницу каталога или архива
		templateFile := "templates/catalog.html"
		if archived {
			templateFile = "templates/archive.html"
		}

		tmpl, err := template.ParseFiles(templateFile)
		if err != nil {
			slog.Error("Ошибка загрузки шаблона", "error", err)
			http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, posts)
	}
}

// HandleCreatePost обрабатывает POST запрос для создания поста
func (h *PostHandler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Получаем данные формы
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		slog.Error("Ошибка парсинга формы", "error", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		name = "Anonymous"
	}

	subject := r.FormValue("subject")
	comment := r.FormValue("comment")

	// Получаем файл изображения
	file, handler, err := r.FormFile("file")
	var imageURL string
	if err == nil && file != nil {
		defer file.Close()

		// TODO: Сохранить файл изображения и получить URL
		// Здесь просто используем заглушку для URL
		slog.Info("Загружен файл", "filename", handler.Filename, "size", handler.Size)
		imageURL = "https://www.google.com/images/branding/googlelogo/2x/googlelogo_light_color_272x92dp.png"
	}

	// Создаем анонимного пользователя
	user, err := h.userService.CreateAnonymousUser(r.Context())
	if err != nil {
		slog.Error("Ошибка создания пользователя", "error", err)
		http.Error(w, "Не удалось создать пользователя", http.StatusInternalServerError)
		return
	}

	// Создаем пост
	post, err := h.postService.CreatePost(r.Context(), subject, comment, imageURL, user.ID)
	if err != nil {
		slog.Error("Ошибка создания поста", "error", err)
		http.Error(w, "Не удалось создать пост", http.StatusInternalServerError)
		return
	}

	// Перенаправляем на страницу созданного поста
	http.Redirect(w, r, "/post/"+strconv.FormatInt(post.ID, 10), http.StatusSeeOther)
}

// HandleArchivePost обрабатывает POST запрос для архивации поста
func (h *PostHandler) HandleArchivePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID поста из пути
	path := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	path = strings.TrimSuffix(path, "/archive")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		slog.Error("Невозможно преобразовать ID в число", "path", path, "error", err)
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	err = h.postService.ArchivePost(r.Context(), id)
	if err != nil {
		slog.Error("Ошибка архивации поста", "id", id, "error", err)
		http.Error(w, "Не удалось архивировать пост", http.StatusInternalServerError)
		return
	}

	// Перенаправляем на страницу архива
	http.Redirect(w, r, "/archive.html", http.StatusSeeOther)
}
