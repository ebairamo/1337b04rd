package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// HandlePage обрабатывает запросы к статическим страницам
func HandlePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request:", r.URL.Path)

	// Обработка формы, если есть данные
	name := r.FormValue("name")
	subject := r.FormValue("subject")
	comment := r.FormValue("comment")
	file, fileHeader, err := r.FormFile("file")

	log.Printf("Form data: name=%s, subject=%s, comment=%s, file=%v, fileHeader=%v, err=%v",
		name, subject, comment, file, fileHeader, err)

	// Определение страницы для отображения
	var page string
	path := r.URL.Path
	if path == "/" {
		page = "templates/catalog.html"
	} else {
		page = "templates" + path
	}

	// Рендеринг шаблона
	tmpl, err := template.ParseFiles(page)
	if err != nil {
		http.Error(w, "Template not found: "+err.Error(), http.StatusNotFound)
		return
	}
	tmpl.Execute(w, nil)
}
