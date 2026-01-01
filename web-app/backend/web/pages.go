package web

import (
	"log"
	"net/http"
	"text/template"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("static/html/homepage.html"))

	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func SingleVideoPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("static/html/singleVideoPage.html"))

	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
		return
	}
}
