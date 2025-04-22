package models

import "time"

// User представляет пользователя в системе
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}
