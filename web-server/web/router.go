package web

import (
	"github.com/gorilla/mux"
)

var RegisterWebRouter = func(router *mux.Router) {
	router.HandleFunc("/", HomePage).Methods("GET")
}
