package services

import (
	"context"
	"log/slog"
	"time"

	"1337b04rd/internal/ports/repositories"
)

// ArchiverService предоставляет функционал для автоматического архивирования постов
type ArchiverService struct {
	postRepo    repositories.PostRepository
	commentRepo repositories.CommentRepository
}

// NewArchiverService создает новый экземпляр сервиса архивирования
func NewArchiverService(postRepo repositories.PostRepository, commentRepo repositories.CommentRepository) *ArchiverService {
	return &ArchiverService{
		postRepo:    postRepo,
		commentRepo: commentRepo,
	}
}

// StartArchiveJob запускает фоновую задачу архивирования
func (s *ArchiverService) StartArchiveJob(ctx context.Context) {
	slog.Info("Запуск фоновой задачи архивирования постов")
	go s.archiveJob(ctx)
}

// archiveJob выполняет периодическую проверку и архивацию постов
func (s *ArchiverService) archiveJob(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute) // Проверка каждую минуту
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.ProcessArchiving(ctx)
		case <-ctx.Done():
			slog.Info("Задача архивирования остановлена")
			return
		}
	}
}

// processArchiving выполняет один цикл проверки и архивирования постов
// Экспортирован для тестирования
func (s *ArchiverService) ProcessArchiving(ctx context.Context) {
	slog.Debug("Запуск проверки постов для архивации")

	// Получаем все неархивированные посты
	posts, err := s.postRepo.GetAllForArchiving(ctx)
	if err != nil {
		slog.Error("Ошибка получения постов для архивации", "error", err)
		return
	}

	now := time.Now()
	archiveCount := 0

	for _, post := range posts {
		// Получаем последний комментарий к посту
		lastComment, err := s.commentRepo.GetLastCommentByPostID(ctx, post.ID)

		if err != nil || lastComment == nil {
			// Пост без комментариев - архивируем через 10 минут
			if now.Sub(post.CreatedAt) > 10*time.Minute {
				err := s.postRepo.Archive(ctx, post.ID)
				if err != nil {
					slog.Error("Ошибка архивирования поста", "post_id", post.ID, "error", err)
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
					continue
				}
				slog.Info("Пост архивирован (15 минут после последнего комментария)",
					"post_id", post.ID,
					"last_comment_at", lastComment.CreatedAt)
				archiveCount++
			}
		}
	}

	if archiveCount > 0 {
		slog.Info("Завершена архивация постов", "archived_count", archiveCount)
	}
}
