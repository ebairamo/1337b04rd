package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"1337b04rd/internal/ports/external"
)

// ImageStorage реализует интерфейс хранилища изображений,
// используя HTTP API S3-хранилища
type ImageStorage struct {
	baseURL    string
	httpClient *http.Client
}

// NewImageStorage создает новый экземпляр хранилища изображений
func NewImageStorage() external.ImageStorage {
	// Получаем хост и порт S3 из переменных окружения или используем значения по умолчанию
	host := os.Getenv("S3_HOST")
	if host == "" {
		host = "localhost" // Используйте localhost для локального тестирования
	}

	port := os.Getenv("S3_PORT")
	if port == "" {
		port = "9000"
	}

	return &ImageStorage{
		baseURL: fmt.Sprintf("http://%s:%s", host, port),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// UploadImage загружает изображение в S3-хранилище
func (s *ImageStorage) UploadImage(ctx context.Context, bucketName, objectKey string, data []byte) (string, error) {
	// Создаем бакет, если он не существует
	if err := s.createBucket(bucketName); err != nil {
		return "", fmt.Errorf("ошибка создания бакета: %w", err)
	}

	// Формируем URL для загрузки объекта
	url := fmt.Sprintf("%s/%s/%s", s.baseURL, bucketName, objectKey)
	slog.Info("Загрузка изображения", "url", url, "size", len(data))

	// Создаем запрос на загрузку
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// Выполняем запрос
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ошибка загрузки изображения: %s, статус: %d", string(body), resp.StatusCode)
	}

	// Возвращаем URL для доступа к изображению
	imageURL := fmt.Sprintf("%s/%s/%s", s.baseURL, bucketName, objectKey)
	slog.Info("Изображение успешно загружено", "bucket", bucketName, "key", objectKey, "url", imageURL)
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

	// Выполняем запрос
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка получения изображения, статус: %d", resp.StatusCode)
	}

	// Читаем данные изображения
	data, err := io.ReadAll(resp.Body)
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

	// Выполняем запрос
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
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

	// Выполняем запрос
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %w", err)
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
	extension := filepath.Ext(originalFilename)
	if extension == "" {
		extension = ".jpg" // Расширение по умолчанию
	}
	return fmt.Sprintf("%d%s", timestamp, extension)
}
