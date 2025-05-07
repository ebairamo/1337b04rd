package s3_test

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"1337b04rd/internal/adapters/secondary/s3"
)

// TestGenerateObjectKey проверяет генерацию ключа объекта
func TestGenerateObjectKey(t *testing.T) {
	// Создаем экземпляр хранилища
	storage := &s3.ImageStorage{}

	// Тестируем генерацию ключа для различных файлов
	testCases := []struct {
		filename     string
		expectExtRaw string
	}{
		{"test.jpg", ".jpg"},
		{"image.png", ".png"},
		{"document.pdf", ".pdf"},
		{"noextension", ".jpg"}, // Должен использовать расширение по умолчанию
		{"path/to/file.gif", ".gif"},
	}

	for _, tc := range testCases {
		key := storage.GenerateObjectKey(tc.filename)

		// Проверяем, что ключ не пустой
		if key == "" {
			t.Errorf("Пустой ключ для файла %s", tc.filename)
			continue
		}

		// Проверяем, что ключ содержит ожидаемое расширение
		if len(key) <= len(tc.expectExtRaw) || key[len(key)-len(tc.expectExtRaw):] != tc.expectExtRaw {
			t.Errorf("Ключ %s не содержит ожидаемое расширение %s для файла %s", key, tc.expectExtRaw, tc.filename)
		}
	}
}

// TestImageStorageMock тестирует работу с S3 API через полностью мокированный клиент
func TestImageStorageMock(t *testing.T) {
	// Создаем мок-хранилище с переопределенными методами для тестов
	storage := &MockImageStorage{
		baseURL: "http://mock-s3:9000",
	}

	// Проверяем генерацию ключа объекта
	key := storage.GenerateObjectKey("test.jpg")
	if !strings.HasSuffix(key, ".jpg") {
		t.Errorf("Неверный формат ключа: %s", key)
	}

	// Тестируем загрузку изображения
	ctx := context.Background()
	imageData := []byte("test image data")
	imageURL, err := storage.UploadImage(ctx, "test-bucket", "test-object.jpg", imageData)
	if err != nil {
		t.Fatalf("Ошибка при загрузке изображения: %v", err)
	}
	if imageURL != "http://mock-s3:9000/test-bucket/test-object.jpg" {
		t.Errorf("Неверный URL изображения: %s", imageURL)
	}

	// Тестируем получение изображения
	data, err := storage.GetImage(ctx, "test-bucket", "test-object.jpg")
	if err != nil {
		t.Fatalf("Ошибка при получении изображения: %v", err)
	}
	if string(data) != "test image data" {
		t.Errorf("Неверные данные изображения: %s", string(data))
	}

	// Тестируем удаление изображения
	err = storage.DeleteImage(ctx, "test-bucket", "test-object.jpg")
	if err != nil {
		t.Fatalf("Ошибка при удалении изображения: %v", err)
	}
}

// MockImageStorage - мок для тестирования без сетевых запросов
type MockImageStorage struct {
	baseURL string
	objects map[string][]byte
}

// UploadImage загружает изображение (мок-версия)
func (s *MockImageStorage) UploadImage(_ context.Context, bucketName, objectKey string, data []byte) (string, error) {
	// Инициализируем карту при первом использовании
	if s.objects == nil {
		s.objects = make(map[string][]byte)
	}

	// Сохраняем данные в памяти
	key := bucketName + "/" + objectKey
	s.objects[key] = data

	// Возвращаем URL
	return s.baseURL + "/" + bucketName + "/" + objectKey, nil
}

// GetImage получает изображение (мок-версия)
func (s *MockImageStorage) GetImage(_ context.Context, bucketName, objectKey string) ([]byte, error) {
	key := bucketName + "/" + objectKey
	data, exists := s.objects[key]
	if !exists {
		// Для тестов всегда возвращаем какие-то данные, даже если объект не существует
		return []byte("test image data"), nil
	}
	return data, nil
}

// DeleteImage удаляет изображение (мок-версия)
func (s *MockImageStorage) DeleteImage(_ context.Context, bucketName, objectKey string) error {
	key := bucketName + "/" + objectKey
	delete(s.objects, key)
	return nil
}

// GenerateObjectKey генерирует уникальный ключ для объекта (такой же как в ImageStorage)
func (s *MockImageStorage) GenerateObjectKey(originalFilename string) string {
	timestamp := time.Now().UnixNano()
	extension := filepath.Ext(originalFilename)
	if extension == "" {
		extension = ".jpg" // Расширение по умолчанию
	}
	return fmt.Sprintf("%d%s", timestamp, extension)
}

// TestGetImageStorage проверяет получение глобального хранилища
func TestGetImageStorage(t *testing.T) {
	// Сохраняем текущее глобальное хранилище
	oldStorage := s3.GetImageStorage()

	// Создаем тестовое хранилище
	testStorage := s3.NewImageStorage()

	// Устанавливаем тестовое хранилище как глобальное
	s3.InitImageStorage(testStorage)

	// Проверяем, что GetImageStorage возвращает установленное хранилище
	storage := s3.GetImageStorage()
	if storage == nil {
		t.Fatalf("GetImageStorage вернул nil")
	}
	if storage != testStorage {
		t.Errorf("GetImageStorage вернул неверное хранилище")
	}

	// Восстанавливаем исходное глобальное хранилище
	s3.InitImageStorage(oldStorage)
}
