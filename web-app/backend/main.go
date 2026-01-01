package main

import (
	"github.com/rs/cors"
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
	socket := shared.WebServerBackendSocket

	router := mux.NewRouter()

	// Adding API router
	api.RegisterAPIRouter(router)

	// printing all routes
	utils.PrintRoutes(router)

	// Adding Middleware - Better do this with a proxy.
	// router.Use(middleware.APIFilter)

	server := &http.Server{
		Handler:      cors.Default().Handler(router),
		Addr:         socket,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Println("Server running on", socket)
	log.Fatal(server.ListenAndServe())
}
