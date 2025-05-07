// internal/adapters/primary/http/middleware/logging.go
// Улучшение логирования для более подробной информации о запросах

package middleware

import (
	"log/slog"
	"net/http"
	"runtime"
	"time"
)

// LoggingMiddleware представляет собой middleware для логирования запросов
type LoggingMiddleware struct {
	DetailedLogging bool // Флаг для включения подробного логирования
}

// NewLoggingMiddleware создает новый экземпляр middleware логирования
func NewLoggingMiddleware(detailed bool) *LoggingMiddleware {
	return &LoggingMiddleware{
		DetailedLogging: detailed,
	}
}

// Handler обрабатывает логирование запросов
func (m *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем ResponseWriter, который может записывать статус
		wrw := newResponseWriter(w)

		// Логируем начало запроса с дополнительной информацией
		// в зависимости от уровня детализации
		if m.DetailedLogging {
			// Более подробное логирование для отладки
			slog.Info("Получен запрос",
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"referer", r.Referer(),
				"content_type", r.Header.Get("Content-Type"),
				"content_length", r.ContentLength,
			)
		} else {
			// Базовое логирование
			slog.Info("Получен запрос",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)
		}

		// Передаем запрос следующему обработчику
		next.ServeHTTP(wrw, r)

		// Логируем окончание запроса
		duration := time.Now().Sub(start)

		// Получаем информацию о системе для мониторинга
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		slog.Info("Запрос обработан",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrw.status,
			"duration_ms", duration.Milliseconds(),
			"bytes_written", wrw.bytesWritten,
			"goroutines", runtime.NumGoroutine(),
			"heap_mb", m.Alloc/1024/1024,
		)
	})
}

// responseWriter это обертка над http.ResponseWriter, которая записывает HTTP статус и количество байт
type responseWriter struct {
	http.ResponseWriter
	status       int
	bytesWritten int64
}

// newResponseWriter создает новый экземпляр responseWriter
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK, 0}
}

// WriteHeader записывает HTTP статус и запоминает его
func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// Write записывает данные и подсчитывает количество байт
func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += int64(n)
	return n, err
}

// Header возвращает заголовки
func (rw *responseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

// Flush реализует интерфейс http.Flusher
func (rw *responseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
