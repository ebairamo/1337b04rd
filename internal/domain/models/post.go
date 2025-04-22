package models

import "time"

// Post представляет пост в системе
type Post struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	ImageURL   string    `json:"image_url,omitempty"`
	UserID     int64     `json:"user_id"`
	UserName   string    `json:"user_name"`
	AvatarURL  string    `json:"avatar_url,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	IsArchived bool      `json:"is_archived"`
}
