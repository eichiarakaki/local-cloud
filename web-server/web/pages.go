package web

import (
	"net/http"
	"text/template"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("static/html/homepage.html"))

	tmpl.Execute(w, nil)
}

func SingleVideoPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("static/html/singleVideoPage.html"))

	tmpl.Execute(w, nil)
}
