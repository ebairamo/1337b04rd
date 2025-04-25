package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// UserHandler обрабатывает запросы к бакетам и объектам
func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		handlePutRequest(w, r)

	case "GET":
		handleGetRequest(w, r)

	case "DELETE":
		handleDeleteRequest(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handlePutRequest обрабатывает PUT запросы
func handlePutRequest(w http.ResponseWriter, r *http.Request) {
	URLS := r.URL.Path[1:]            // Убираем начальный слэш
	parts := strings.Split(URLS, "/") // Разделяем по "/"

	if len(parts) == 1 {
		// Создаем бакет
		BucketName := parts[0]

		// Проверяем имя бакета на валидность
		if valid, errResponse := checkBucketName(BucketName); !valid {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(errResponse.Code)

			xmlResponse, err := xml.Marshal(errResponse)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Write(xmlResponse)
			return
		}

		// Проверяем, существует ли уже бакет с таким именем

		if _, err := os.Stat(BaseDir + "/" + BucketName); err == nil {
			// Если бакет существует, возвращаем ошибку 409 Conflict
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusConflict)

			errorResponse := ErrorResponse{
				Code:    http.StatusConflict,
				Message: "Bucket already exists",
			}

			xmlResponse, err := xml.Marshal(errorResponse)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Write(xmlResponse)
		}

		// Создаем бакет

		if err := os.MkdirAll(BaseDir+"/"+BucketName, 0o777); err != nil {

			http.Error(w, "Failed to create bucket", http.StatusInternalServerError)
			return
		}

		// Добавляем информацию о создании бакета
		CreationTime := time.Now().Format(time.RFC3339)
		Status := "marked for deletion"
		addBucketToCSV(BucketName, CreationTime, Status)

		w.WriteHeader(http.StatusOK)
		return
	}

	if len(parts) == 2 {
		BucketName, ObjectName := parts[0], parts[1]
		if valid, errResponse := checkBucketName(BucketName); !valid {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(errResponse.Code)

			xmlResponse, err := xml.Marshal(errResponse)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Write(xmlResponse)
			return
		}
		pathToBucket := filepath.Join(BaseDir, BucketName)
		// Проверяем, существует ли бакет

		if _, err := os.Stat(pathToBucket); os.IsNotExist(err) {
			http.Error(w, "Bucket does not exist", http.StatusNotFound)
			return
		}

		// Получаем текущее время для LastModified
		LastModified := time.Now().Format(time.RFC3339)
		ObjectDir := filepath.Join(pathToBucket, ObjectName)

		// Создаем файл для сохранения объекта
		if err := saveObjectFile(r.Body, ObjectDir); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		validBucketName := BucketName

		// Получаем размер файла
		size, err := getFileSize(ObjectDir)
		if err != nil {
			http.Error(w, "Failed to get file info", http.StatusInternalServerError)
			return
		}

		// Записываем данные в CSV
		if err := addObjectToCSV(ObjectName, size, LastModified, pathToBucket, validBucketName); err != nil {
			http.Error(w, "Failed to update CSV", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/xml")
		return
	}

	http.Error(w, "Invalid request", http.StatusBadRequest)
}

// handleGetRequest обрабатывает GET запросы
func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	URLS := r.URL.Path[1:]            // Убираем начальный слэш
	parts := strings.Split(URLS, "/") // Разделяем по "/"

	if len(URLS) == 0 {
		// Список бакетов
		buckets, err := listBuckets()
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Failed to list buckets")
			return
		}

		response := BucketList{Buckets: buckets}
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(response)
		return
	}
	if len(parts) == 1 {

		object, err := listObject(string(parts[0]))
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Failed to list buckets")

			return
		}

		response := ObjectList{Objects: object}
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(response)
		return
	}

	if len(parts) == 2 {
		// Получение конкретного объекта из бакета
		BucketName := parts[0]
		ObjectName := parts[1]
		objectPath := filepath.Join(BaseDir, BucketName, ObjectName)

		// Чтение файла (например, изображения)
		content, err := os.ReadFile(objectPath)
		if err != nil {
			if os.IsNotExist(err) {
				ErrorHandler(w, http.StatusNotFound, "Object not found")
				return
			}
			ErrorHandler(w, http.StatusInternalServerError, "Failed to read object")
			return
		}
		contentType, err := getFileContentType(objectPath)
		if err != nil {
		}
		// Установка заголовков для отправки изображения
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
		w.WriteHeader(http.StatusOK)
		w.Write(content)
		return
	}

	http.Error(w, "Invalid request", http.StatusBadRequest)
}

// handleDeleteRequest обрабатывает DELETE запросы
func handleDeleteRequest(w http.ResponseWriter, r *http.Request) {
	URLS := r.URL.Path[1:]
	parts := strings.Split(URLS, "/")

	if len(parts) < 1 {
		http.Error(w, "Bucket name required", http.StatusBadRequest)
		return
	}

	BucketName := parts[0]
	// projectPath, _ := os.Getwd()
	pathToBucket := filepath.Join(BaseDir, BucketName)
	var ObjectName string
	// if len(parts) == 1 {
	// }
	if len(parts) == 1 {
		// Проверяем, существует ли корзина
		if _, err := os.Stat(pathToBucket); os.IsNotExist(err) {
			// Корзина не существует
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusNotFound)

			errorResponse := ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Bucket not found",
			}
			xmlResponse, err := xml.Marshal(errorResponse)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Write(xmlResponse)
			return
		}

		// Проверяем, пустая ли корзина или содержит только objects.csv
		entries, err := os.ReadDir(pathToBucket)
		if err != nil {
			http.Error(w, "Error reading directory: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Если корзина содержит файлы, кроме objects.csv
		if len(entries) > 1 || (len(entries) == 1 && entries[0].Name() != "objects.csv") {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusConflict)

			errorResponse := ErrorResponse{
				Code:    http.StatusConflict,
				Message: "Bucket is not empty. Cannot delete.",
			}
			xmlResponse, err := xml.Marshal(errorResponse)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Write(xmlResponse)
			return
		}

		// Если корзина пуста или содержит только objects.csv, удаляем ее
		err = os.RemoveAll(pathToBucket)
		if err != nil {
			http.Error(w, "Error removing bucket: "+err.Error(), http.StatusInternalServerError)
			return
		} else {
			// Обновляем статус в CSV
			updateCSVStatus("buckets.csv", BucketName)
			// Возвращаем код состояния 204 No Content
			w.WriteHeader(http.StatusNoContent)
		}
	}

	if len(parts) == 2 {
		ObjectName = parts[1]
		if ObjectName != "" {
			// Проверяем, существует ли объект в метаданных CSV перед удалением
			filePath := filepath.Join(pathToBucket, "objects.csv")
			file, err := os.Open(filePath)
			if err != nil {
				http.Error(w, "Error opening metadata file", http.StatusInternalServerError)
				return
			}
			defer file.Close()

			reader := csv.NewReader(file)
			records, err := reader.ReadAll()
			if err != nil {
				http.Error(w, "Error reading metadata file", http.StatusInternalServerError)
				return
			}

			objectExists := false
			for _, record := range records {
				if record[0] == ObjectName { // Предполагаем, что objectName находится в первой колонке
					objectExists = true
					break
				}
			}

			if !objectExists {
				http.Error(w, "Object not found", http.StatusNotFound)
				return
			}

			// Удаляем файлы и обновляем CSV, если объект найден
			if err := removeAllFilesInDir(pathToBucket, ObjectName); err != nil {
				http.Error(w, "Error removing files", http.StatusInternalServerError)
				return
			}

			if err := removeObjectFromCSV(pathToBucket, ObjectName); err != nil {
				http.Error(w, "Error updating object metadata", http.StatusInternalServerError)
				return
			}
			updateLastModifiedInBucketCSV(BucketName)

			// Возвращаем статус 204 No Content после успешного удаления
			w.WriteHeader(http.StatusNoContent)
			fmt.Println("Object and associated files removed successfully.")
			return
		}
	}
}

// Функция для удаления объекта из CSV
func removeObjectFromCSV(bucketName, objectName string) error {
	filePath := filepath.Join(bucketName, "objects.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	var updatedRows [][]string
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV records: %w", err)
	}

	for _, record := range records {
		if record[0] != objectName { // Предполагаем, что objectName находится в первой колонке
			updatedRows = append(updatedRows, record)
		}
	}

	file, err = os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(updatedRows); err != nil {
		return fmt.Errorf("failed to write updated CSV records: %w", err)
	}

	return nil
}

// saveObjectFile сохраняет данные объекта в файл
func saveObjectFile(body io.ReadCloser, path string) error {
	defer body.Close() // Закрываем r.Body
	content, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}
	if len(content) == 0 {
		return fmt.Errorf("empty file")
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

// getFileSize возвращает размер файла
func getFileSize(path string) (string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}
	return strconv.FormatInt(fileInfo.Size(), 10), nil
}
