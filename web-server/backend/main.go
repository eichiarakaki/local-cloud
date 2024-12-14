package main

import (
	"log"
	"net/http"
	shared "shared_mods"
	"time"

	"github.com/eichiarakaki/local-cloud/web-server/backend/api"
	"github.com/eichiarakaki/local-cloud/web-server/backend/utils"
	"github.com/gorilla/mux"
)

func main() {
	err := shared.LoadConfig("../../config.json")
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	socket := shared.WebServerSocket

	router := mux.NewRouter()
	// Adding Web router
	//web.RegisterWebRouter(router) // No need anymore because frontend is in a other folder.

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
