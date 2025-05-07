// internal/adapters/primary/http/templates.go

package http

import (
	"html/template"
	"net/http"
	"path/filepath"
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
func RenderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}, title string, pageTitle string) error {
	// Находим все шаблоны
	templatePath := filepath.Join("templates", templateName)
	baseTemplatePath := filepath.Join("templates", "base.html")

	// Парсим шаблоны
	tmpl, err := template.ParseFiles(baseTemplatePath, templatePath)
	if err != nil {
		return err
	}

	// Подготавливаем данные для шаблона
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
