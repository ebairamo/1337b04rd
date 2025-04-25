package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

func listObject(filePath string) ([]Object, error) {
	var objects []Object
	filePath = BaseDir + filePath
	bucketFilePath := filepath.Join(filePath, "objects.csv")
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
		size := int64(0)
		Object := Object{
			ObjectKey:    record[0],
			Size:         size,
			ContentType:  record[2],
			LastModified: record[3],
		}
		objects = append(objects, Object)
	}

	return objects, nil
}

func updateCSVStatusObject(filePath string, ObjectName string) error {
	baseDir := filePath
	objects, err := listObject(baseDir)
	if err != nil {
		return err
	}

	objectFilePath := filepath.Join(baseDir, "Object.csv")
	file, err := os.OpenFile(objectFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Запись заголовка, если файл пустой
	fileInfo, err := file.Stat()
	if err == nil && fileInfo.Size() == 0 {
		if err := writer.Write([]string{"Name", "Size", "ContentType", "LastModified"}); err != nil {
			return fmt.Errorf("error writing header to CSV: %v", err)
		}
	}

	// Перебор объектов и запись только тех, что не совпадают с ObjectName
	for _, object := range objects {
		if object.ObjectKey != ObjectName {
			if err := writer.Write([]string{object.ObjectKey, fmt.Sprintf("%d", object.Size), object.ContentType, object.LastModified}); err != nil {
				return fmt.Errorf("error writing object %s to CSV: %v", object.ObjectKey, err)
			}
		}
	}

	return nil
}
