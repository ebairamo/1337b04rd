package services

import (
	"context"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"

	"1337b04rd/internal/ports/repositories"
)

// ArchiverService предоставляет функционал для автоматического архивирования постов
type ArchiverService struct {
	postRepo      repositories.PostRepository
	commentRepo   repositories.CommentRepository
	interval      time.Duration
	lastRun       time.Time
	statsLock     sync.Mutex
	archivedCount int
	errorCount    int
	isRunning     bool
}

// ArchiverStats содержит статистику работы архиватора
type ArchiverStats struct {
	LastRun       time.Time
	ArchivedCount int
	ErrorCount    int
	IsRunning     bool
}

// NewArchiverService создает новый экземпляр сервиса архивирования
func NewArchiverService(postRepo repositories.PostRepository, commentRepo repositories.CommentRepository) *ArchiverService {
	return &ArchiverService{
		postRepo:    postRepo,
		commentRepo: commentRepo,
		interval:    1 * time.Minute, // По умолчанию проверка каждую минуту
	}
}

// SetInterval устанавливает интервал проверки архивирования
func (s *ArchiverService) SetInterval(interval time.Duration) {
	s.interval = interval
}

// StartArchiveJob запускает фоновую задачу архивирования
func (s *ArchiverService) StartArchiveJob(ctx context.Context) {
	slog.Info("Запуск фоновой задачи архивирования постов", "interval", s.interval)
	s.isRunning = true
	go s.archiveJob(ctx)
}

// GetStats возвращает статистику работы архиватора
func (s *ArchiverService) GetStats() ArchiverStats {
	s.statsLock.Lock()
	defer s.statsLock.Unlock()

	return ArchiverStats{
		LastRun:       s.lastRun,
		ArchivedCount: s.archivedCount,
		ErrorCount:    s.errorCount,
		IsRunning:     s.isRunning,
	}
}

// archiveJob выполняет периодическую проверку и архивацию постов
func (s *ArchiverService) archiveJob(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Используем recover для предотвращения остановки горутины при панике
			func() {
				defer func() {
					if r := recover(); r != nil {
						stack := debug.Stack()
						slog.Error("Паника в задаче архивирования",
							"panic", r,
							"stack", string(stack),
						)
						// Увеличиваем счетчик ошибок
						s.statsLock.Lock()
						s.errorCount++
						s.statsLock.Unlock()
					}
				}()

				s.ProcessArchiving(ctx)
			}()
		case <-ctx.Done():
			s.isRunning = false
			slog.Info("Задача архивирования остановлена")
			return
		}
	}
}

// ProcessArchiving выполняет один цикл проверки и архивирования постов
func (s *ArchiverService) ProcessArchiving(ctx context.Context) {
	// Обновляем время последнего запуска
	s.statsLock.Lock()
	s.lastRun = time.Now()
	s.statsLock.Unlock()

	slog.Debug("Запуск проверки постов для архивации")

	// Получаем все неархивированные посты
	posts, err := s.postRepo.GetAllForArchiving(ctx)
	if err != nil {
		slog.Error("Ошибка получения постов для архивации", "error", err)
		// Увеличиваем счетчик ошибок
		s.statsLock.Lock()
		s.errorCount++
		s.statsLock.Unlock()
		return
	}

	now := time.Now()
	archiveCount := 0

	for _, post := range posts {
		// Получаем последний комментарий к посту
		lastComment, err := s.commentRepo.GetLastCommentByPostID(ctx, post.ID)

		if err != nil {
			slog.Error("Ошибка получения последнего комментария", "post_id", post.ID, "error", err)
			// Увеличиваем счетчик ошибок
			s.statsLock.Lock()
			s.errorCount++
			s.statsLock.Unlock()
			continue
		}

		if lastComment == nil {
			// Пост без комментариев - архивируем через 10 минут
			if now.Sub(post.CreatedAt) > 10*time.Minute {
				err := s.postRepo.Archive(ctx, post.ID)
				if err != nil {
					slog.Error("Ошибка архивирования поста", "post_id", post.ID, "error", err)
					// Увеличиваем счетчик ошибок
					s.statsLock.Lock()
					s.errorCount++
					s.statsLock.Unlock()
					continue
				}
				slog.Info("Пост архивирован (без комментариев)", "post_id", post.ID, "created_at", post.CreatedAt)
				archiveCount++
			}
		} else {
			// Пост с комментариями - архивируем через 15 минут после последнего комментария
			if now.Sub(lastComment.CreatedAt) > 15*time.Minute {
				err := s.postRepo.Archive(ctx, post.ID)
				if err != nil {
					slog.Error("Ошибка архивирования поста", "post_id", post.ID, "error", err)
					// Увеличиваем счетчик ошибок
					s.statsLock.Lock()
					s.errorCount++
					s.statsLock.Unlock()
					continue
				}
				slog.Info("Пост архивирован (15 минут после последнего комментария)",
					"post_id", post.ID,
					"last_comment_at", lastComment.CreatedAt)
				archiveCount++
			}
		}
	}

	// Обновляем статистику
	if archiveCount > 0 {
		s.statsLock.Lock()
		s.archivedCount += archiveCount
		s.statsLock.Unlock()
		slog.Info("Завершена архивация постов", "archived_count", archiveCount, "total_archived", s.archivedCount)
	}
}
