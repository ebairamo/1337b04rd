package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

func listBuckets() ([]Bucket, error) {
	var buckets []Bucket

	bucketFilePath := filepath.Join(BaseDir, "buckets.csv")
	file, err := os.Open(bucketFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read() // Skip header
	if err != nil {
		return nil, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		bucket := Bucket{
			Name:             record[0],
			CreationTime:     record[1],
			LastModifiedTime: record[2],
			Status:           record[3],
		}
		buckets = append(buckets, bucket)
	}

	return buckets, nil
}

func updateCSVStatus(filePath string, bucketName string) error {
	baseDir := BaseDir
	buckets, err := listBuckets()
	if err != nil {
		return err
	}

	bucketFilePath := filepath.Join(baseDir, filePath)
	file, err := os.OpenFile(bucketFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Name", "CreationTime", "LastModifiedTime", "Status"}) // Запись заголовка

	// Перебор корзин и запись только тех, что не совпадают с bucketName
	for _, bucket := range buckets {
		if bucket.Name != bucketName {
			writer.Write([]string{bucket.Name, bucket.CreationTime, bucket.LastModifiedTime, bucket.Status})
		}
	}

	return nil
}

func removeAllFilesInDir(dirPath string, ObjectName string) error {
	// Читаем содержимое директории
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}

	// Удаляем каждый файл или папку
	for _, entry := range entries {
		if entry.Name() == ObjectName {
			fullPath := filepath.Join(dirPath, entry.Name()) // Полный путь
			if err := os.RemoveAll(fullPath); err != nil {
				return fmt.Errorf("error removing %s: %v", fullPath, err)
			}
		}
	}
	updateCSVStatusObject(dirPath, ObjectName)
	return nil
}
