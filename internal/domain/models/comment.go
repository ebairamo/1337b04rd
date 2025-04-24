package models

import "time"

// Comment представляет комментарий в системе
type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	UserName  string    `json:"user_name"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	ReplyToID int64     `json:"reply_to_id,omitempty"`
}
