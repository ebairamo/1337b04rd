// internal/adapters/secondary/s3/image_storage.go

package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"1337b04rd/internal/ports/external"
)

const (
	// Максимальный размер файла в байтах (5 МБ)
	MaxFileSize = 5 * 1024 * 1024

	// Разрешенные типы изображений
	AllowedImageTypes = "image/jpeg,image/png,image/gif,image/svg+xml,image/webp"
)

// ImageStorageOptions содержит настройки для хранилища изображений
type ImageStorageOptions struct {
	MaxFileSize       int64
	AllowedImageTypes []string
}

// ImageStorage реализует интерфейс хранилища изображений,
// используя HTTP API S3-хранилища
type ImageStorage struct {
	baseURL    string
	httpClient *http.Client
	options    ImageStorageOptions
}

// NewImageStorage создает новый экземпляр хранилища изображений
func NewImageStorage() external.ImageStorage {
	// Получаем хост и порт S3 из переменных окружения или используем значения по умолчанию
	host := os.Getenv("S3_HOST")
	if host == "" {
		host = "localhost" // Используем localhost для локального тестирования
	}

	port := os.Getenv("S3_PORT")
	if port == "" {
		port = "9000"
	}

	// Настройки по умолчанию
	options := ImageStorageOptions{
		MaxFileSize:       MaxFileSize,
		AllowedImageTypes: strings.Split(AllowedImageTypes, ","),
	}

	return &ImageStorage{
		baseURL: fmt.Sprintf("http://%s:%s", host, port),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		options: options,
	}
}

// UploadImage загружает изображение в S3-хранилище с проверками
func (s *ImageStorage) UploadImage(ctx context.Context, bucketName, objectKey string, data []byte) (string, error) {
	// Проверяем размер файла
	if len(data) > int(s.options.MaxFileSize) {
		return "", fmt.Errorf("размер файла превышает допустимый предел %d байт", s.options.MaxFileSize)
	}

	// Определяем тип файла
	fileType := http.DetectContentType(data)

	// Проверяем, является ли тип файла разрешенным изображением
	allowed := false
	for _, allowedType := range s.options.AllowedImageTypes {
		if fileType == allowedType {
			allowed = true
			break
		}
	}

	if !allowed {
		slog.Warn("Попытка загрузки файла неподдерживаемого типа", "type", fileType)
		return "", fmt.Errorf("неподдерживаемый тип файла: %s. Разрешены только изображения", fileType)
	}

	// Создаем бакет, если он не существует
	if err := s.createBucket(bucketName); err != nil {
		return "", fmt.Errorf("ошибка создания бакета: %w", err)
	}

	// Формируем URL для загрузки объекта
	url := fmt.Sprintf("%s/%s/%s", s.baseURL, bucketName, objectKey)
	slog.Info("Загрузка изображения", "url", url, "size", len(data), "type", fileType)

	// Создаем запрос на загрузку
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// Устанавливаем Content-Type
	req.Header.Set("Content-Type", fileType)

	// Выполняем запрос с повторными попытками при ошибках соединения
	var resp *http.Response
	maxRetries := 3

	for retry := 0; retry < maxRetries; retry++ {
		resp, err = s.httpClient.Do(req)

		if err != nil {
			// Если это не последняя попытка, попробуем еще раз
			if retry < maxRetries-1 {
				slog.Warn("Ошибка соединения при загрузке, повторная попытка", "retry", retry+1, "error", err)
				time.Sleep(time.Duration(retry+1) * 500 * time.Millisecond) // Экспоненциальная задержка
				continue
			}
			return "", fmt.Errorf("ошибка выполнения запроса после %d попыток: %w", maxRetries, err)
		}
		break // Если успешно, выходим из цикла
	}

	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ошибка загрузки изображения: %s, статус: %d", string(body), resp.StatusCode)
	}

	// Возвращаем URL для доступа к изображению
	imageURL := fmt.Sprintf("%s/%s/%s", s.baseURL, bucketName, objectKey)
	slog.Info("Изображение успешно загружено", "bucket", bucketName, "key", objectKey, "url", imageURL, "type", fileType)
	return imageURL, nil
}

// GetImage получает изображение из S3-хранилища
func (s *ImageStorage) GetImage(ctx context.Context, bucketName, objectKey string) ([]byte, error) {
	// Формируем URL для получения объекта
	url := fmt.Sprintf("%s/%s/%s", s.baseURL, bucketName, objectKey)

	// Создаем запрос на получение
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// Выполняем запрос с повторными попытками
	var resp *http.Response
	maxRetries := 3

	for retry := 0; retry < maxRetries; retry++ {
		resp, err = s.httpClient.Do(req)

		if err != nil {
			// Если это не последняя попытка, попробуем еще раз
			if retry < maxRetries-1 {
				slog.Warn("Ошибка соединения при получении изображения, повторная попытка", "retry", retry+1, "error", err)
				time.Sleep(time.Duration(retry+1) * 500 * time.Millisecond)
				continue
			}
			return nil, fmt.Errorf("ошибка выполнения запроса после %d попыток: %w", maxRetries, err)
		}
		break
	}

	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, errors.New("изображение не найдено")
		}
		return nil, fmt.Errorf("ошибка получения изображения, статус: %d", resp.StatusCode)
	}

	// Проверяем Content-Type
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		slog.Warn("Получен объект с неожиданным Content-Type", "content_type", contentType)
	}

	// Читаем данные изображения с ограничением размера
	var maxSize int64 = MaxFileSize
	if contentLength := resp.ContentLength; contentLength > 0 && contentLength > maxSize {
		return nil, fmt.Errorf("размер изображения (%d байт) превышает максимально допустимый (%d байт)", contentLength, maxSize)
	}

	reader := io.LimitReader(resp.Body, maxSize)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения данных изображения: %w", err)
	}

	return data, nil
}

// DeleteImage удаляет изображение из S3-хранилища
func (s *ImageStorage) DeleteImage(ctx context.Context, bucketName, objectKey string) error {
	// Формируем URL для удаления объекта
	url := fmt.Sprintf("%s/%s/%s", s.baseURL, bucketName, objectKey)

	// Создаем запрос на удаление
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// Выполняем запрос с повторными попытками
	var resp *http.Response
	maxRetries := 3

	for retry := 0; retry < maxRetries; retry++ {
		resp, err = s.httpClient.Do(req)

		if err != nil {
			if retry < maxRetries-1 {
				slog.Warn("Ошибка соединения при удалении изображения, повторная попытка", "retry", retry+1, "error", err)
				time.Sleep(time.Duration(retry+1) * 500 * time.Millisecond)
				continue
			}
			return fmt.Errorf("ошибка выполнения запроса после %d попыток: %w", maxRetries, err)
		}
		break
	}

	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			// Объект не найден, но мы считаем, что это не ошибка при удалении
			slog.Info("Изображение не найдено при попытке удаления", "bucket", bucketName, "key", objectKey)
			return nil
		}
		return fmt.Errorf("ошибка удаления изображения, статус: %d", resp.StatusCode)
	}

	slog.Info("Изображение успешно удалено", "bucket", bucketName, "key", objectKey)
	return nil
}

// createBucket создает бакет в S3-хранилище
func (s *ImageStorage) createBucket(bucketName string) error {
	// Формируем URL для создания бакета
	url := fmt.Sprintf("%s/%s", s.baseURL, bucketName)
	slog.Info("Создание бакета", "url", url)

	// Создаем запрос на создание бакета
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// Выполняем запрос с повторными попытками
	var resp *http.Response
	maxRetries := 3

	for retry := 0; retry < maxRetries; retry++ {
		resp, err = s.httpClient.Do(req)

		if err != nil {
			if retry < maxRetries-1 {
				slog.Warn("Ошибка соединения при создании бакета, повторная попытка", "retry", retry+1, "error", err)
				time.Sleep(time.Duration(retry+1) * 500 * time.Millisecond)
				continue
			}
			return fmt.Errorf("ошибка выполнения запроса после %d попыток: %w", maxRetries, err)
		}
		break
	}

	defer resp.Body.Close()

	// Статус 200 или 409 (уже существует) считаем успехом
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusConflict {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ошибка создания бакета: %s, статус: %d", string(body), resp.StatusCode)
	}

	slog.Info("Бакет успешно создан или уже существует", "bucket", bucketName)
	return nil
}

// GenerateObjectKey генерирует уникальный ключ для объекта
func (s *ImageStorage) GenerateObjectKey(originalFilename string) string {
	timestamp := time.Now().UnixNano()

	// Получаем расширение файла
	extension := filepath.Ext(originalFilename)
	if extension == "" {
		// Пытаемся определить расширение из имени файла
		dotIndex := strings.LastIndex(originalFilename, ".")
		if dotIndex != -1 {
			extension = originalFilename[dotIndex:]
		} else {
			// Расширение по умолчанию
			extension = ".jpg"
		}
	}

	// Нормализуем расширение, приводя к нижнему регистру
	extension = strings.ToLower(extension)

	// Проверяем расширение на соответствие MIME-типу изображения
	mimeType := mime.TypeByExtension(extension)
	if mimeType == "" || !strings.HasPrefix(mimeType, "image/") {
		// Если расширение не соответствует изображению, используем .jpg
		extension = ".jpg"
	}

	return fmt.Sprintf("%d%s", timestamp, extension)
}

// ValidateImageData проверяет данные изображения
func (s *ImageStorage) ValidateImageData(data []byte) error {
	// Проверяем размер файла
	if len(data) > int(s.options.MaxFileSize) {
		return fmt.Errorf("размер файла превышает допустимый предел %d байт", s.options.MaxFileSize)
	}

	// Проверяем тип файла
	fileType := http.DetectContentType(data)

	// Проверяем, является ли тип файла разрешенным изображением
	allowed := false
	for _, allowedType := range s.options.AllowedImageTypes {
		if fileType == allowedType {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("неподдерживаемый тип файла: %s. Разрешены только изображения", fileType)
	}

	return nil
}
