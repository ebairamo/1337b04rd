package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func HandlePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request:", r.URL.Path)

	// Обработка формы, если есть данные
	name := r.FormValue("name")
	subject := r.FormValue("subject")
	comment := r.FormValue("comment")
	file, fileHeader, err := r.FormFile("file")

	fmt.Println("file:", file, "fileHeader:", fileHeader, "err:", err,
		"name:", name, "subject:", subject, "comment:", comment)

	// Определение страницы для отображения
	var page string
	path := r.URL.Path
	if path == "/" {
		page = "templates/catalog.html"
	} else {
		page = "templates" + path
	}

	// Рендеринг шаблона
	tmpl, _ := template.ParseFiles(page)
	tmpl.Execute(w, nil)
}
