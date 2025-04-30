package http

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"1337b04rd/internal/adapters/primary/http/handlers"
	"1337b04rd/internal/adapters/primary/http/middleware"
	"1337b04rd/internal/adapters/secondary/postgres"
	"1337b04rd/internal/adapters/secondary/rickandmorty"
	"1337b04rd/internal/adapters/secondary/s3"
	"1337b04rd/internal/domain/services"
)

// RegisterRoutes регистрирует все маршруты приложения
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	// Инициализация сервисов и репозиториев
	avatarService := rickandmorty.NewAvatarService()
	imageStorage := s3.NewImageStorage()
	userRepo := postgres.NewUserRepository(db, avatarService)
	postRepo := postgres.NewPostRepository(db)
	commentRepo := postgres.NewCommentRepository(db)

	userService := services.NewUserService(userRepo)
	postService := services.NewPostService(postRepo, userRepo)
	commentService := services.NewCommentService(commentRepo, userRepo, postRepo)

	// Инициализируем глобальное хранилище для использования в обработчиках
	s3.InitImageStorage(imageStorage)

	// Создание middleware
	authMiddleware := middleware.NewAuthMiddleware(userService)
	loggingMiddleware := middleware.NewLoggingMiddleware()

	// Создание обработчиков
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService, userService, commentService)
	commentHandler := handlers.NewCommentHandler(commentService, userService)
	pageHandler := handlers.HandlePage

	// Функция-помощник для оборачивания обработчиков с аутентификацией
	withAuth := func(handler http.Handler) http.Handler {
		return loggingMiddleware.Handler(authMiddleware.Handler(handler))
	}

	// Регистрация маршрутов для API с аутентификацией
	mux.Handle("/api/", withAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Маршруты для пользователей
		if strings.HasPrefix(path, "/api/users/") {
			handleUserRoutes(w, r, userHandler)
			return
		}

		// Маршруты для постов
		if strings.HasPrefix(path, "/api/posts/") {
			handlePostRoutes(w, r, postHandler, commentHandler)
			return
		}

		// Маршруты для комментариев
		if strings.HasPrefix(path, "/api/comments/") {
			handleCommentRoutes(w, r, commentHandler)
			return
		}

		http.NotFound(w, r)
	})))

	// Маршруты для работы с постами
	mux.Handle("/post/", withAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/post/")
		r.URL.Path = "/api/posts/" + id
		postHandler.HandleGetPost(w, r)
	})))

	// Маршруты для отправки форм
	mux.Handle("/submit-post", withAuth(http.HandlerFunc(postHandler.HandleCreatePost)))

	// Обработчик для создания комментариев
	mux.Handle("/submit-comment", withAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Получен запрос на создание комментария: %s %s", r.Method, r.URL.Path)
		commentHandler.HandleCreateComment(w, r)
	})))

	// Маршруты для страниц каталога и архива
	mux.Handle("/catalog.html", withAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postHandler.HandleGetAllPosts(w, r)
	})))

	mux.Handle("/archive.html", withAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		q.Set("archived", "true")
		r.URL.RawQuery = q.Encode()
		postHandler.HandleGetAllPosts(w, r)
	})))

	// Статические страницы без аутентификации
	mux.Handle("/", withAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// Перенаправляем на каталог при запросе корневой страницы
			http.Redirect(w, r, "/catalog.html", http.StatusFound)
			return
		}

		// Для всех остальных запросов используем обработчик страниц
		pageHandler(w, r)
	})))
}

// handleUserRoutes обрабатывает маршруты пользователей
func handleUserRoutes(w http.ResponseWriter, r *http.Request, handler *handlers.UserHandler) {
	switch r.Method {
	case http.MethodGet:
		handler.HandleGetUser(w, r)
	case http.MethodPost:
		handler.HandleCreateUser(w, r)
	default:
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
	}
}

// handlePostRoutes обрабатывает маршруты постов
func handlePostRoutes(w http.ResponseWriter, r *http.Request, postHandler *handlers.PostHandler, commentHandler *handlers.CommentHandler) {
	path := r.URL.Path

	// Маршрут для архивации поста
	if strings.HasSuffix(path, "/archive") && r.Method == http.MethodPost {
		postHandler.HandleArchivePost(w, r)
		return
	}

	// Проверка на маршрут комментариев
	if strings.HasSuffix(path, "/comments") {
		switch r.Method {
		case http.MethodGet:
			commentHandler.HandleGetPostComments(w, r)
		default:
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
		return
	}

	// Маршрут для получения всех постов
	if path == "/api/posts/" {
		switch r.Method {
		case http.MethodGet:
			postHandler.HandleGetAllPosts(w, r)
		default:
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
		return
	}

	// Маршрут для получения поста по ID
	switch r.Method {
	case http.MethodGet:
		postHandler.HandleGetPost(w, r)
	default:
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
	}
}

// handleCommentRoutes обрабатывает маршруты комментариев
func handleCommentRoutes(w http.ResponseWriter, r *http.Request, handler *handlers.CommentHandler) {
	switch r.Method {
	case http.MethodGet:
		handler.HandleGetComment(w, r)
	case http.MethodDelete:
		handler.HandleDeleteComment(w, r)
	default:
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
	}
}
