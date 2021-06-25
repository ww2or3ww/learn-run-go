package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("views/templates/%s.html", file))
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}

func StartMainServer() error {
	files := http.FileServer(http.Dir("views"))
	http.Handle("/static/", http.StripPrefix("/static/", files))

	http.HandleFunc("/", hello)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/world", world)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	log.Println(fmt.Sprintf("port=%v", port))

	return http.ListenAndServe(":"+port, nil)
}
