package handlers

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

// TemplateData содержит общие данные для всех шаблонов
type TemplateData struct {
	Title       string
	PageTitle   string
	CurrentYear int
	Data        interface{}
}

// RenderTemplate рендерит шаблон с данными
func RenderTemplate(w http.ResponseWriter, templateName string, data interface{}, title string, pageTitle string) error {
	// Загружаем шаблоны
	tmpl, err := template.ParseFiles("templates/base.html", "templates/"+templateName)
	if err != nil {
		return err
	}

	// Создаем данные для шаблона, включая текущий год
	templateData := TemplateData{
		Title:       title,
		PageTitle:   pageTitle,
		CurrentYear: time.Now().Year(),
		Data:        data,
	}

	// Рендерим шаблон
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.Execute(w, templateData)
}

// HandlePage обрабатывает запросы к статическим страницам
func HandlePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос:", r.URL.Path)

	// Определение страницы для отображения
	var templateName string
	var title string
	var pageTitle string

	path := r.URL.Path
	if path == "/" {
		templateName = "catalog.html"
		title = "Каталог"
		pageTitle = "Каталог постов"
	} else {
		templateName = path[1:] // Убираем начальный слэш
		switch templateName {
		case "create-post.html":
			title = "Создать пост"
			pageTitle = "Создание нового поста"
		case "archive.html":
			title = "Архив"
			pageTitle = "Архив постов"
		default:
			title = "1337b04rd"
			pageTitle = "Доска сообщений"
		}
	}

	// Рендеринг шаблона
	err := RenderTemplate(w, templateName, nil, title, pageTitle)
	if err != nil {
		log.Printf("Ошибка при рендеринге шаблона %s: %s", templateName, err)
		http.Error(w, "Шаблон не найден: "+err.Error(), http.StatusNotFound)
		return
	}
}
