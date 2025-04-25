package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	BaseDir string
	Port    string // Порт по умолчанию
	Help    bool
)

func init() {
	flag.StringVar(&BaseDir, "dir", "data", "DIRECTORY")
	flag.StringVar(&Port, "port", "8080", "PORT")
	flag.BoolVar(&Help, "help", false, "HELP")
	flag.Parse()
}

func main() {
	// Проверяем наличие аргумента для вывода справки
	if Help {
		printHelp()
		return
	}
	projectPath, _ := os.Getwd()
	if IsDirPathInside(projectPath, BaseDir) {
		initBaseDir()
	} else {
		fmt.Println("Файл находится вне проекта")
		os.Exit(1)
	}
	// projectPath, _ := os.Getwd()
	// print(projectPath)
	// currentDir := filepath.Join(projectPath, "data")
	// if !IsDirPathInside(projectPath, currentDir) {
	// 	fmt.Println("Invalid Project name or path")
	// 	os.Exit(1)
	// }

	// Инициализируем базовую директорию, если указана
	// dir := getDirFromArgs(args)
	// if dir != "" {
	// 	initBaseDir(dir)
	// }
	// initBaseDir()

	http.Handle("/", http.HandlerFunc(UserHandler))
	err := http.ListenAndServe(":"+Port, nil)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		return
	}
	fmt.Printf("Сервер работает на порту %s...\n", Port)
}

// func IsDirPathInside(projectPath, dirPath string) bool {
// 	absProjectPath, _ := filepath.Abs(projectPath)
// 	absDirPath, _ := filepath.Abs(dirPath)

// 	return absProjectPath == filepath.Dir(absDirPath) || strings.HasPrefix(absDirPath, absProjectPath+string(os.PathListSeparator))
// }

// Инициализация базовой директории для хранения
func initBaseDir() {
	// print(BaseDir)
	if err := os.MkdirAll(BaseDir, os.ModePerm); err != nil {
		fmt.Println("Не удалось создать базовую директорию:", err)
		os.Exit(1)
	}
}

// Получаем директорию из аргументов

func IsDirPathInside(projectPath, dirPath string) bool {
	absProjectPath, _ := filepath.Abs(projectPath)
	absDirPath, _ := filepath.Abs(dirPath)

	return absProjectPath == filepath.Dir(absDirPath) || strings.HasPrefix(absDirPath, absProjectPath+string(os.PathListSeparator))
}
