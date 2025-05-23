package handlers

import (
	"encoding/json"
	"html/template"
	"io"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"1337b04rd/internal/adapters/primary/http/middleware"
	"1337b04rd/internal/adapters/secondary/s3"
	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/domain/services"
)

// PostHandler обрабатывает HTTP запросы для постов
type PostHandler struct {
	postService    *services.PostService
	userService    *services.UserService
	commentService *services.CommentService
}

// NewPostHandler создает новый обработчик постов
func NewPostHandler(postService *services.PostService, userService *services.UserService, commentService *services.CommentService) *PostHandler {
	return &PostHandler{
		postService:    postService,
		userService:    userService,
		commentService: commentService,
	}
}

// PaginationData содержит информацию о пагинации для шаблонов
type PaginationData struct {
	Posts       []*models.Post
	CurrentPage int
	PrevPage    int
	NextPage    int
	TotalPages  int
	PageNumbers []int
	Limit       int
}

// FixImageURL преобразует URL изображения для доступа из браузера

func FixImageURL(url string) string {
	// Заменяем http://s3:9000/ на /s3-proxy/
	return strings.Replace(url, "http://s3:9000", "http://localhost:9000", 1)
}

// HandleGetPost обрабатывает GET запрос для получения поста
func (h *PostHandler) HandleGetPost(w http.ResponseWriter, r *http.Request) {
	// Получаем пользователя из контекста
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		slog.Error("Пользователь не найден в контексте")
		http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
		return
	}

	// Простой парсинг пути для извлечения ID поста
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

	// Исправляем URL изображения для доступа из браузера
	if post.ImageURL != "" {
		post.ImageURL = FixImageURL(post.ImageURL)
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

		// Получаем комментарии к посту
		comments, err := h.commentService.GetCommentsByPostID(r.Context(), post.ID, 50, 0)
		if err != nil {
			slog.Error("Ошибка при получении комментариев", "post_id", post.ID, "error", err)
			http.Error(w, "Ошибка при получении комментариев", http.StatusInternalServerError)
			return
		}

		// Исправляем URL изображений в комментариях
		for i := range comments {
			if comments[i].ImageURL != "" {
				comments[i].ImageURL = FixImageURL(comments[i].ImageURL)
			}
		}

		// Создаем данные для шаблона
		templateData := struct {
			*models.Post
			Comments []*models.Comment
			User     *models.User
		}{
			Post:     post,
			Comments: comments,
			User:     user,
		}

		// Передаем данные в шаблон
		tmpl.Execute(w, templateData)
	}
}

// HandleGetAllPosts обрабатывает GET запрос для получения списка постов с пагинацией
func (h *PostHandler) HandleGetAllPosts(w http.ResponseWriter, r *http.Request) {
	// Получаем пользователя из контекста
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		slog.Error("Пользователь не найден в контексте")
		http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
		return
	}

	// Параметры запроса
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	archivedStr := r.URL.Query().Get("archived")

	// Устанавливаем значения по умолчанию и парсим параметры
	page := 1
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	limit := 10 // По умолчанию 10 постов на страницу
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := (page - 1) * limit

	archived := false
	if archivedStr == "true" || archivedStr == "1" {
		archived = true
	}

	// Получаем общее количество постов для расчета пагинации
	posts, err := h.postService.GetAllPosts(r.Context(), limit, offset, archived)
	if err != nil {
		slog.Error("Ошибка получения списка постов", "error", err)
		http.Error(w, "Не удалось получить список постов", http.StatusInternalServerError)
		return
	}

	// Исправляем URL изображений для всех постов
	for i := range posts {
		if posts[i].ImageURL != "" {
			posts[i].ImageURL = FixImageURL(posts[i].ImageURL)
		}
	}
	// Исправляем URL изображений для всех постов
	for i := range posts {
		if posts[i].ImageURL != "" {
			posts[i].ImageURL = FixImageURL(posts[i].ImageURL)
		}
	}

	// Используем предполагаемое общее количество постов для демонстрации
	// В реальном приложении следует получить точное количество из базы данных
	totalPosts := 100 // Предполагаемое количество для демонстрации
	if archived {
		totalPosts = 50 // Меньше постов в архиве
	}

	totalPages := int(math.Ceil(float64(totalPosts) / float64(limit)))
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}
	nextPage := page + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	// Создаем список номеров страниц для отображения в пагинации
	pageNumbers := make([]int, 0)
	startPage := math.Max(1, float64(page-2))
	endPage := math.Min(float64(totalPages), float64(page+2))

	for i := startPage; i <= endPage; i++ {
		pageNumbers = append(pageNumbers, int(i))
	}

	// Определяем значения заголовков в зависимости от того, архив это или каталог
	var title, pageTitle string
	if archived {
		title = "Архив"
		pageTitle = "Архив постов"
	} else {
		title = "Каталог"
		pageTitle = "Каталог постов"
	}

	// Теперь создаем данные для шаблона
	templateData := struct {
		Title       string
		PageTitle   string
		CurrentYear int
		Posts       []*models.Post
		CurrentPage int
		PrevPage    int
		NextPage    int
		TotalPages  int
		PageNumbers []int
		Limit       int
	}{
		Title:       title,
		PageTitle:   pageTitle,
		CurrentYear: time.Now().Year(),
		Posts:       posts,
		CurrentPage: page,
		PrevPage:    prevPage,
		NextPage:    nextPage,
		TotalPages:  totalPages,
		PageNumbers: pageNumbers,
		Limit:       limit,
	}

	// Определяем, какой шаблон использовать
	templateName := "catalog.html"
	if archived {
		templateName = "archive.html"
	}

	// Сначала парсим base.html, а затем конкретный шаблон
	tmpl, err := template.ParseFiles("templates/base.html", "templates/"+templateName)
	if err != nil {
		slog.Error("Ошибка загрузки шаблона", "template", templateName, "error", err)
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}

	// Передаем данные в шаблон
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, templateData)
	if err != nil {
		slog.Error("Ошибка выполнения шаблона", "template", templateName, "error", err)
		http.Error(w, "Ошибка рендеринга", http.StatusInternalServerError)
		return
	}
}

// HandleCreatePost обрабатывает POST запрос для создания поста
func (h *PostHandler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	// Получаем данные формы
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		slog.Error("Ошибка парсинга формы", "error", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Получаем имя пользователя из формы или используем имя из сессии
	name := r.FormValue("name")
	if name != "" && name != user.Username {
		// В реальной реализации здесь нужно обновить имя пользователя в БД
		user.Username = name
	}

	subject := r.FormValue("subject")
	comment := r.FormValue("comment")

	// Получаем файл изображения
	file, handler, err := r.FormFile("file")
	var imageURL string
	if err == nil && file != nil {
		defer file.Close()

		// Читаем содержимое файла
		buffer, err := io.ReadAll(file)
		if err != nil {
			slog.Error("Ошибка чтения файла", "error", err)
			http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
			return
		}

		// Генерируем ключ для объекта и загружаем изображение, если доступно хранилище
		if storage := s3.GetImageStorage(); storage != nil {
			objectKey := storage.GenerateObjectKey(handler.Filename)
			imageURL, err = storage.UploadImage(r.Context(), "posts", objectKey, buffer)
			if err != nil {
				slog.Error("Ошибка загрузки изображения", "error", err)
				// Продолжаем без изображения, если произошла ошибка
				imageURL = ""
			} else {
				slog.Info("Изображение загружено", "filename", handler.Filename, "size", handler.Size, "url", imageURL)
			}
		} else {
			slog.Warn("Хранилище изображений не инициализировано")
		}
	}

	// Создаем пост, используя ID пользователя из сессии
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

	// Получаем пользователя из контекста
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		slog.Error("Пользователь не найден в контексте")
		http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
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
