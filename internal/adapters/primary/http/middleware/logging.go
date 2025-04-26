package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// LoggingMiddleware представляет собой middleware для логирования запросов
type LoggingMiddleware struct{}

// NewLoggingMiddleware создает новый экземпляр middleware логирования
func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

// Handler обрабатывает логирование запросов
func (m *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем ResponseWriter, который может записывать статус
		wrw := newResponseWriter(w)

		// Логируем начало запроса
		slog.Info("Получен запрос",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		// Передаем запрос следующему обработчику
		next.ServeHTTP(wrw, r)

		// Логируем окончание запроса
		duration := time.Since(start)
		slog.Info("Запрос обработан",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrw.status,
			"duration", duration,
		)
	})
}

// responseWriter это обертка над http.ResponseWriter, которая записывает HTTP статус
type responseWriter struct {
	http.ResponseWriter
	status int
}

// newResponseWriter создает новый экземпляр responseWriter
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

// WriteHeader записывает HTTP статус и запоминает его
func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}
