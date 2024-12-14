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

	// Videos storage
	// videoStoragePath, err := filepath.Abs("../videos-storage/")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Serving videos from directory:", videoStoragePath) // Debug
	// APIRouter.PathPrefix("/videos-storage/").
	// Handler(http.StripPrefix("/videos-storage/", http.FileServer(http.Dir(videoStoragePath)))).
	// Methods("GET")
	APIRouter.HandleFunc("/videos-storage/{mediaName}", ServeStorage).Methods("GET")
}
