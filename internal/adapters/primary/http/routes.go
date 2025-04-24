package http

import (
	"1337b04rd/internal/adapters/primary/http/handlers"
	"1337b04rd/internal/adapters/secondary/postgres"
	"1337b04rd/internal/domain/services"
	"database/sql"
	"net/http"
	"strings"
)

// RegisterRoutes регистрирует все маршруты приложения
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	// Базовый обработчик для статических страниц
	mux.HandleFunc("/", handlers.HandlePage)

	// Регистрация маршрутов для пользователей, постов и комментариев
	RegisterUserRoutes(mux, db)
	RegisterPostRoutes(mux, db)
	RegisterCommentRoutes(mux, db)
}

// RegisterUserRoutes регистрирует маршруты для пользователей
func RegisterUserRoutes(mux *http.ServeMux, db *sql.DB) {
	// Инициализация репозитория и сервиса
	userRepo := postgres.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Регистрация маршрутов
	mux.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.HandleGetUser(w, r)
		case http.MethodPost:
			userHandler.HandleCreateUser(w, r)
		default:
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
	})
}

// RegisterPostRoutes регистрирует маршруты для постов
func RegisterPostRoutes(mux *http.ServeMux, db *sql.DB) {
	// Инициализация репозиториев и сервисов
	postRepo := postgres.NewPostRepository(db)
	userRepo := postgres.NewUserRepository(db)

	userService := services.NewUserService(userRepo)
	postService := services.NewPostService(postRepo, userRepo)

	postHandler := handlers.NewPostHandler(postService, userService)

	// Обработка API маршрутов
	mux.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Маршрут для архивации поста
		if strings.HasSuffix(path, "/archive") && r.Method == http.MethodPost {
			postHandler.HandleArchivePost(w, r)
			return
		}

		// Проверка на маршрут комментариев (будет обработан в RegisterCommentRoutes)
		if strings.HasSuffix(path, "/comments") {
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
	})

	// Обработка маршрута создания поста
	mux.HandleFunc("/submit-post", postHandler.HandleCreatePost)

	// Обработка маршрутов шаблонов
	mux.HandleFunc("/catalog.html", func(w http.ResponseWriter, r *http.Request) {
		postHandler.HandleGetAllPosts(w, r)
	})

	mux.HandleFunc("/archive.html", func(w http.ResponseWriter, r *http.Request) {
		// Установка параметра archived=true для запроса архивных постов
		q := r.URL.Query()
		q.Set("archived", "true")
		r.URL.RawQuery = q.Encode()

		postHandler.HandleGetAllPosts(w, r)
	})

	// Обработка маршрута для отдельного поста
	mux.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/post/")
		r.URL.Path = "/api/posts/" + id

		postHandler.HandleGetPost(w, r)
	})
}

// RegisterCommentRoutes регистрирует маршруты для комментариев
func RegisterCommentRoutes(mux *http.ServeMux, db *sql.DB) {
	// Инициализация репозиториев и сервисов
	commentRepo := postgres.NewCommentRepository(db)
	userRepo := postgres.NewUserRepository(db)
	postRepo := postgres.NewPostRepository(db)

	userService := services.NewUserService(userRepo)
	commentService := services.NewCommentService(commentRepo, userRepo, postRepo)

	commentHandler := handlers.NewCommentHandler(commentService, userService)

	// Регистрация маршрутов для работы с отдельными комментариями
	mux.HandleFunc("/api/comments/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			commentHandler.HandleGetComment(w, r)
		case http.MethodDelete:
			commentHandler.HandleDeleteComment(w, r)
		default:
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
	})

	// Регистрация маршрутов для комментариев к постам
	mux.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/comments") {
			switch r.Method {
			case http.MethodGet:
				commentHandler.HandleGetPostComments(w, r)
			default:
				http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			}
		}
	})

	// Обработка отправки комментария
	mux.HandleFunc("/submit-comment", commentHandler.HandleCreateComment)
}
