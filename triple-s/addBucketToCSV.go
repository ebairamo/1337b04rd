package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Функция для добавления или обновления информации о корзине в CSV
func addBucketToCSV(bucketName string, creationTime string, status string) error {
	csvFilePath := filepath.Join(BaseDir, "buckets.csv")
	fileExists := false
	if _, err := os.Stat(csvFilePath); err == nil {
		fileExists = true
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

	// Проверяем существующие записи
	for i, record := range records {
		if record[0] == bucketName {
			// Обновляем запись о корзине
			records[i][2] = creationTime // LastModifiedTime - обновляем его тоже на время создания
			records[i][3] = status       // Обновляем статус
			updated = true
			break
		}
	}

	// Если запись не найдена, добавляем новую
	if !updated {
		records = append(records, []string{bucketName, creationTime, creationTime, status})
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

	return nil
}

func updateLastModifiedInBucketCSV(bucketName string) error {
	// Путь к файлу buckets.csv
	filePath := filepath.Join(BaseDir, "buckets.csv")

	// Открываем файл для чтения
	file, err := os.OpenFile(filePath, os.O_RDWR, 0o644)
	if err != nil {
		return fmt.Errorf("error opening buckets metadata file: %w", err)
	}
	defer file.Close()

	// Чтение всех записей из файла
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading buckets metadata file: %w", err)
	}

	// Обновляем время последней модификации для нужного bucket
	lastModified := time.Now().Format(time.RFC3339)
	bucketPath := filepath.Join(BaseDir, bucketName)
	isEmpty := false
	isEmpty, _ = isBucketEmpty(bucketPath)

	for i, record := range records {
		if record[0] == bucketName {
			if isEmpty {
				record[3] = "Marked for deletion"
			}
			// Предполагаем, что имя bucket находится в первой колонке
			// Обновляем время последней модификации (например, в колонке с индексом 2)

			// if err != nil {
			// 	fmt.Println("Ошибка:", err)
			if len(record) > 2 {
				record[2] = lastModified
			} else {
				record = append(record, lastModified)
			}
			records[i] = record
			break
		}
	}

	// Перезаписываем файл с обновленными записями
	file.Seek(0, 0)
	file.Truncate(0)
	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("error writing updated metadata to buckets file: %w", err)
	}

	return nil
}

func isBucketEmpty(bucketPath string) (bool, error) {
	dir, _ := os.Open(bucketPath)
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return false, err
	}
	return files[0].Name() == "objects.csv", nil

	// entries, err := os.ReadDir(bucketPath)
	// if err != nil {
	// 	return false, err // Ошибка при чтении директории
	// }
	// return len(entries) == 0, nil // true, если нет файлов
}
