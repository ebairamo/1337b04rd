package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func getFileContentType(filePath string) (string, error) {
	// Проверка расширения файла
	ext := filepath.Ext(filePath)
	switch ext {
	case ".txt":
		return "text/plain", nil
	case ".jpg", ".jpeg":
		return "image/jpeg", nil
	case ".png":
		return "image/png", nil
	// Добавьте другие расширения по мере необходимости
	default:
		// Если расширение неизвестно, попытайтесь определить тип содержимого
		return detectContentTypeFromFile(filePath)
	}
}

func detectContentTypeFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	return http.DetectContentType(buf), nil
}

// Функция для добавления или обновления объекта в CSV
// Функция для добавления или обновления объекта в CSV
func addObjectToCSV(objectName string, size string, lastModified string, pathToBucket string, validBucketName string) error {
	// Создаем директорию для бакета, если она не существует
	if err := os.MkdirAll(pathToBucket, os.ModePerm); err != nil {
		return fmt.Errorf("error creating bucket directory: %v", err)
	}

	// Полный путь к файлу
	filePath := filepath.Join(pathToBucket, objectName)

	// Определяем MIME-тип файла
	contentType, err := getFileContentType(filePath)
	if err != nil {
		return fmt.Errorf("error getting content type: %v", err)
	}

	// Полный путь к CSV
	csvFilePath := filepath.Join(pathToBucket, "objects.csv")

	// Проверяем, существует ли файл
	fileExists := false
	if _, err := os.Stat(csvFilePath); err == nil {
		fileExists = true
	}

	// Открываем файл для добавления данных
	file, err := os.OpenFile(csvFilePath, os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV: %v", err)
	}

	// Проверяем, существует ли запись для данного объекта
	found := false
	for i, record := range records {
		if record[0] == objectName {
			// Обновляем только LastModified, Size и ContentType
			records[i][1] = size
			records[i][2] = contentType
			records[i][3] = lastModified
			found = true
			break
		}
	}

	// Если запись не найдена, добавляем новую
	if !found {
		records = append(records, []string{objectName, size, contentType, lastModified})
	}

	// Открываем файл для записи и очищаем его
	file.Close() // Закрываем файл для записи
	file, err = os.OpenFile(csvFilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		return fmt.Errorf("error opening file for writing: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Если файл новый, добавляем заголовки
	if !fileExists {
		headers := []string{"ObjectName", "Size", "ContentType", "LastModified"}
		if err := writer.Write(headers); err != nil {
			return fmt.Errorf("error writing headers to CSV: %v", err)
		}
	}

	// Записываем все записи обратно в файл
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing to CSV: %v", err)
		}
	}
	updateBucketToCSV(validBucketName)
	return nil
}

func updateBucketToCSV(bucketName string) error {
	csvFilePath := filepath.Join(BaseDir, "buckets.csv")

	// Создаем директорию для хранения CSV, если она не существует
	if err := os.MkdirAll(BaseDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating data directory: %v", err)
	}

	// Открываем файл для чтения и записи
	file, err := os.OpenFile(csvFilePath, os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV: %v", err)
	}

	// Флаг для отслеживания, была ли обновлена запись
	updated := false
	lastModifiedTime := time.Now().Format(time.RFC3339)
	status := "Active"
	// Проверяем существующие записи
	for i, record := range records {
		if record[0] == bucketName {
			// Обновляем только LastModifiedTime
			records[i][3] = status
			records[i][2] = lastModifiedTime // LastModifiedTime
			updated = true
			break
		}
	}

	// Если запись не найдена, добавляем новую с начальным временем создания и последней модификацией одинаковыми
	if !updated {
		creationTime := lastModifiedTime // Время создания равно времени модификации
		records = append(records, []string{bucketName, creationTime, lastModifiedTime, status})
	}

	// Открываем файл для записи и очищаем его
	file.Close() // Закрываем файл для записи
	file, err = os.OpenFile(csvFilePath, os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		return fmt.Errorf("error opening file for writing: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Если файл новый, добавляем заголовки
	if !updated {
		headers := []string{"Name", "CreationTime", "LastModifiedTime", "Status"}
		if err := writer.Write(headers); err != nil {
			return fmt.Errorf("error writing headers to CSV: %v", err)
		}
	}

	// Записываем все записи обратно в файл
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing to CSV: %v", err)
		}
	}

	return nil
}
