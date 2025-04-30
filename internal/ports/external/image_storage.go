package external

import "context"

// ImageStorage представляет интерфейс для работы с хранилищем изображений
type ImageStorage interface {
	// UploadImage загружает изображение в хранилище и возвращает URL для доступа к нему
	UploadImage(ctx context.Context, bucketName, objectKey string, data []byte) (string, error)

	// GetImage получает изображение из хранилища по имени бакета и ключу объекта
	GetImage(ctx context.Context, bucketName, objectKey string) ([]byte, error)

	// DeleteImage удаляет изображение из хранилища
	DeleteImage(ctx context.Context, bucketName, objectKey string) error

	// GenerateObjectKey генерирует уникальный ключ для объекта
	GenerateObjectKey(originalFilename string) string
}
