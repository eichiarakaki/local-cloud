package api

import (
	"github.com/gorilla/mux"
)

var RegisterAPIRouter = func(router *mux.Router) {
	APIRouter := router.PathPrefix("/api").Subrouter()

	APIRouter.HandleFunc("/videos", GetAllVideos).Methods("GET")
	APIRouter.HandleFunc("/video/{videoID}", GetVideoByID).Methods("GET")

	// Just for testing routes
	APIRouter.HandleFunc("/test/{id}", Testing).Methods("GET")

	APIRouter.HandleFunc("/videos-storage/{mediaName}", ServeStorage).Methods("GET")

	APIRouter.HandleFunc("/send-to-downloader-server", SendToDownloaderServer).Methods("POST")
}
