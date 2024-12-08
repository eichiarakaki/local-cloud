package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

var RegisterWebRouter = func(router *mux.Router) {
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	router.HandleFunc("/", HomePage).Methods("GET")
}
