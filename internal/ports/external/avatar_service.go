package external

import "context"

// AvatarService представляет интерфейс для получения аватаров
type AvatarService interface {
	// GetRandomAvatar возвращает URL и имя случайного аватара
	GetRandomAvatar(ctx context.Context) (string, string, error)

	// ResetUsedIDs сбрасывает список использованных ID
	ResetUsedIDs()
}
