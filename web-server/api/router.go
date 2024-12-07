package api

import (
	"github.com/gorilla/mux"
)

var RegisterAPIRouter = func(router *mux.Router) {
	APIRouter := router.PathPrefix("/api").Subrouter()

	APIRouter.HandleFunc("/video/{videoID}", GetVideoByID).Methods("GET")
}