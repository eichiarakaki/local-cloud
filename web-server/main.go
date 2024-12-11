package main

import (
	"log"
	"net/http"
	shared "shared_mods"
	"time"

	"github.com/eichiarakaki/local-cloud/api"
	"github.com/eichiarakaki/local-cloud/utils"
	"github.com/eichiarakaki/local-cloud/web"
	"github.com/gorilla/mux"
)

func main() {
	err := shared.LoadConfig()
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	socket := shared.WebServerSocket

	router := mux.NewRouter()
	// Adding Web router
	web.RegisterWebRouter(router)
	// Adding API router
	api.RegisterAPIRouter(router)

	// printing all routes
	utils.PrintRoutes(router)

	// Adding Middleware - Better do this with a proxy.
	// router.Use(middleware.APIFilter)

	server := &http.Server{
		Handler:      router,
		Addr:         socket,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server running on", socket)
	log.Fatal(server.ListenAndServe())
}
