package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"1337b04rd/internal/domain/models"
)

// CommentRepository реализует интерфейс репозитория комментариев для PostgreSQL
type CommentRepository struct {
	db *sql.DB
}

// NewCommentRepository создает новый экземпляр репозитория комментариев
func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

// GetByID возвращает комментарий по его ID
func (r *CommentRepository) GetByID(ctx context.Context, id int64) (*models.Comment, error) {
	slog.Info("Заглушка: получение комментария", "id", id)
	return &models.Comment{
		ID:        id,
		PostID:    1,
		UserID:    1,
		UserName:  "Anonymous",
		AvatarURL: "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
		Content:   "Тестовый комментарий",
		CreatedAt: time.Now(),
	}, nil
}

// GetByPostID возвращает все комментарии к указанному посту
func (r *CommentRepository) GetByPostID(ctx context.Context, postID int64, limit, offset int) ([]*models.Comment, error) {
	// Устанавливаем значения по умолчанию для параметров пагинации
	if limit <= 0 {
		limit = 50 // Значение по умолчанию
	}
	if offset < 0 {
		offset = 0
	}

	// SQL-запрос с выборкой всех полей
	query := `SELECT 
        id, 
        post_id, 
        user_id, 
        user_name, 
        avatar_url, 
        content, 
        image_url, 
        created_at, 
        reply_to_id
        FROM comments 
        WHERE post_id = $1 
        ORDER BY created_at ASC 
        LIMIT $2 OFFSET $3`

	// Выполняем запрос
	slog.Info("Выполнение запроса комментариев",
		"post_id", postID,
		"limit", limit,
		"offset", offset)

	rows, err := r.db.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		slog.Error("Ошибка выполнения запроса комментариев",
			"post_id", postID,
			"error", err.Error())
		return nil, fmt.Errorf("ошибка запроса комментариев: %w", err)
	}
	defer rows.Close()

	// Создаем слайс для результатов
	var comments []*models.Comment

	// Итерируемся по результатам
	for rows.Next() {
		var comment models.Comment

		// Временная переменная для обработки NULL-значений
		var replyToID sql.NullInt64
		var imageURL sql.NullString
		var avatarURL sql.NullString

		// Сканируем строку в структуру, обрабатывая возможные NULL-значения
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.UserName,
			&avatarURL,
			&comment.Content,
			&imageURL,
			&comment.CreatedAt,
			&replyToID,
		)
		if err != nil {
			slog.Error("Ошибка сканирования строки комментария",
				"error", err.Error())
			return nil, fmt.Errorf("ошибка сканирования: %w", err)
		}

		// Присваиваем значения из Nullable-полей
		if avatarURL.Valid {
			comment.AvatarURL = avatarURL.String
		}
		if imageURL.Valid {
			comment.ImageURL = imageURL.String
		}
		if replyToID.Valid {
			comment.ReplyToID = replyToID.Int64
		}

		// Добавляем комментарий в результаты
		comments = append(comments, &comment)
	}

	// Проверяем наличие ошибок итерации
	if err = rows.Err(); err != nil {
		slog.Error("Ошибка после итерации по комментариям",
			"error", err.Error())
		return nil, fmt.Errorf("ошибка итерации: %w", err)
	}

	// Логируем количество найденных комментариев
	slog.Info("Найдены комментарии",
		"post_id", postID,
		"count", len(comments))

	return comments, nil
}

// Create создает новый комментарий
func (r *CommentRepository) Create(ctx context.Context, comment *models.Comment) (int64, error) {
	slog.Info("Заглушка: создание комментария",
		"post_id", comment.PostID,
		"user_id", comment.UserID,
		"content", comment.Content)
	return 1, nil
}

// Delete удаляет комментарий по ID
func (r *CommentRepository) Delete(ctx context.Context, id int64) error {
	slog.Info("Заглушка: удаление комментария", "id", id)
	return nil
}
