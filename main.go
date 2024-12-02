package main

import (
	"log"
	"net/http"
	"time"

	"github.com/eichiarakaki/local-cloud/api"
	"github.com/eichiarakaki/local-cloud/middleware"
	"github.com/eichiarakaki/local-cloud/web"
	"github.com/gorilla/mux"
)

func main() {
	socket := "localhost:3030"
	router := mux.NewRouter()
	// Adding Web router
	web.RegisterWebRouter(router)
	// Adding API router
	api.RegisterAPIRouter(router)

	// Adding Middleware
	router.Use(middleware.APIFilter)

	server := &http.Server{
		Handler:      router,
		Addr:         socket,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server running on", socket)
	log.Fatal(server.ListenAndServe())
}
