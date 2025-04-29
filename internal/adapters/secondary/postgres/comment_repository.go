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
	query := `SELECT 
        id, post_id, user_id, user_name, avatar_url, content, image_url, created_at, reply_to_id
        FROM comments 
        WHERE id = $1`

	var comment models.Comment
	var avatarURL, imageURL sql.NullString
	var replyToID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
		if err == sql.ErrNoRows {
			slog.Error("Комментарий не найден", "id", id)
			return nil, fmt.Errorf("комментарий с id %d не найден", id)
		}
		slog.Error("Ошибка получения комментария", "id", id, "error", err)
		return nil, err
	}

	// Обрабатываем null-значения
	if avatarURL.Valid {
		comment.AvatarURL = avatarURL.String
	}
	if imageURL.Valid {
		comment.ImageURL = imageURL.String
	}
	if replyToID.Valid {
		comment.ReplyToID = replyToID.Int64
	}

	slog.Info("Комментарий получен", "id", comment.ID, "post_id", comment.PostID)
	return &comment, nil
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
        id, post_id, user_id, user_name, avatar_url, content, image_url, created_at, reply_to_id
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
	// SQL запрос на вставку комментария
	query := `INSERT INTO comments 
        (post_id, user_id, user_name, avatar_url, content, image_url, created_at, reply_to_id) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
        RETURNING id`

	slog.Info("Создание комментария",
		"post_id", comment.PostID,
		"user_id", comment.UserID,
		"content_length", len(comment.Content))

	var id int64
	var replyToID interface{} = nil

	// Устанавливаем replyToID, если он не равен 0
	if comment.ReplyToID > 0 {
		replyToID = comment.ReplyToID
	}

	// Выполняем запрос
	err := r.db.QueryRowContext(ctx, query,
		comment.PostID,
		comment.UserID,
		comment.UserName,
		comment.AvatarURL,
		comment.Content,
		comment.ImageURL,
		time.Now(),
		replyToID,
	).Scan(&id)
	if err != nil {
		slog.Error("Ошибка при создании комментария", "error", err.Error())
		return 0, fmt.Errorf("ошибка создания комментария: %w", err)
	}

	slog.Info("Комментарий успешно создан", "id", id, "post_id", comment.PostID)
	return id, nil
}

// Delete удаляет комментарий по ID
func (r *CommentRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM comments WHERE id = $1`

	slog.Info("Удаление комментария", "id", id)

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.Error("Ошибка при удалении комментария", "id", id, "error", err.Error())
		return fmt.Errorf("ошибка удаления комментария: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Ошибка получения количества затронутых строк", "error", err.Error())
		return fmt.Errorf("ошибка получения количества затронутых строк: %w", err)
	}

	if rowsAffected == 0 {
		slog.Warn("Комментарий не найден", "id", id)
		return fmt.Errorf("комментарий с id %d не найден", id)
	}

	slog.Info("Комментарий успешно удален", "id", id)
	return nil
}
