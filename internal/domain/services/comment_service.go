package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/ports/repositories"
)

// CommentService предоставляет бизнес-логику для работы с комментариями
type CommentService struct {
	commentRepo repositories.CommentRepository
	userRepo    repositories.UserRepository
	postRepo    repositories.PostRepository
}

// NewCommentService создает новый экземпляр сервиса комментариев
func NewCommentService(
	commentRepo repositories.CommentRepository,
	userRepo repositories.UserRepository,
	postRepo repositories.PostRepository,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		userRepo:    userRepo,
		postRepo:    postRepo,
	}
}

// GetCommentByID возвращает комментарий по ID
func (s *CommentService) GetCommentByID(ctx context.Context, id int64) (*models.Comment, error) {
	slog.Info("Получение комментария по ID", "id", id)
	return s.commentRepo.GetByID(ctx, id)
}

// GetCommentsByPostID возвращает комментарии к посту
func (s *CommentService) GetCommentsByPostID(ctx context.Context, postID int64, limit, offset int) ([]*models.Comment, error) {
	slog.Info("Получение комментариев к посту", "post_id", postID, "limit", limit, "offset", offset)
	return s.commentRepo.GetByPostID(ctx, postID, limit, offset)
}

// CreateComment создает новый комментарий
func (s *CommentService) CreateComment(
	ctx context.Context,
	postID int64,
	userID int64,
	content string,
	imageURL string,
	replyToID int64,
) (*models.Comment, error) {
	// Проверяем существование поста
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		slog.Error("Ошибка при получении поста", "post_id", postID, "error", err)
		return nil, fmt.Errorf("пост не найден: %w", err)
	}

	if post.IsArchived {
		slog.Warn("Попытка создать комментарий к архивному посту", "post_id", postID, "user_id", userID)
		return nil, fmt.Errorf("нельзя комментировать архивные посты")
	}

	// Получаем информацию о пользователе
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		slog.Error("Ошибка при получении пользователя", "user_id", userID, "error", err)
		// В случае ошибки, используем дефолтные значения
		user = &models.User{
			ID:        userID,
			Username:  "Anonymous",
			AvatarURL: "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
		}
	}

	// Проверяем существование родительского комментария, если указан
	if replyToID > 0 {
		_, err := s.commentRepo.GetByID(ctx, replyToID)
		if err != nil {
			slog.Error("Ошибка при получении родительского комментария", "reply_to_id", replyToID, "error", err)
			// Если родительский комментарий не найден, сбрасываем replyToID
			replyToID = 0
		}
	}

	// Создаем объект комментария
	comment := &models.Comment{
		PostID:    postID,
		UserID:    userID,
		UserName:  user.Username,
		AvatarURL: user.AvatarURL,
		Content:   content,
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
		ReplyToID: replyToID,
	}

	// Сохраняем комментарий в БД
	commentID, err := s.commentRepo.Create(ctx, comment)
	if err != nil {
		slog.Error("Ошибка при сохранении комментария", "error", err)
		return nil, fmt.Errorf("не удалось сохранить комментарий: %w", err)
	}

	// Устанавливаем ID созданного комментария
	comment.ID = commentID

	slog.Info("Комментарий успешно создан",
		"comment_id", comment.ID,
		"post_id", postID,
		"user_id", userID,
		"reply_to_id", replyToID)

	return comment, nil
}

// DeleteComment удаляет комментарий
func (s *CommentService) DeleteComment(ctx context.Context, id int64) error {
	slog.Info("Удаление комментария", "id", id)
	return s.commentRepo.Delete(ctx, id)
}
