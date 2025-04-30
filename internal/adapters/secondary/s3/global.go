package s3

import (
	"1337b04rd/internal/ports/external"
)

// Глобальное хранилище изображений для использования в обработчиках
var GlobalImageStorage external.ImageStorage

// InitImageStorage инициализирует глобальное хранилище изображений
func InitImageStorage(storage external.ImageStorage) {
	GlobalImageStorage = storage
}

// GetImageStorage возвращает глобальное хранилище изображений
func GetImageStorage() external.ImageStorage {
	return GlobalImageStorage
}
