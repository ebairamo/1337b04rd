package services

import (
	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/ports/repositories"
	"context"
	"log/slog"
	"time"
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
	slog.Info("Получение комментариев к посту", "post_id", postID)
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
	// В заглушке просто возвращаем статические данные
	slog.Info("Создание комментария",
		"post_id", postID,
		"user_id", userID,
		"content_length", len(content))

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		// Если не удалось получить пользователя, используем анонима
		user = &models.User{
			ID:        userID,
			Username:  "Anonymous",
			AvatarURL: "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
		}
	}

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

	// В реальной реализации здесь было бы сохранение в БД
	comment.ID = 1 // Заглушка для ID

	return comment, nil
}

// DeleteComment удаляет комментарий
func (s *CommentService) DeleteComment(ctx context.Context, id int64) error {
	slog.Info("Удаление комментария", "id", id)
	return s.commentRepo.Delete(ctx, id)
}
