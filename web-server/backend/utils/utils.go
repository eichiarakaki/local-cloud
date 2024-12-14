package utils

import (
	"fmt"
	"log"

	"github.com/gorilla/mux"
)

func PrintRoutes(r *mux.Router) {
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			// If there is no template, the route probably does not have a Path defined
			t = "<no template>"
		}

		// Get methods associated with the route
		methods, err := route.GetMethods()
		if err != nil {
			// If there are no methods, we leave an empty list
			methods = []string{"<no methods>"}
		}

		fmt.Printf("Route: %s, Methods: %v\n", t, methods)
		return nil
	})

	if err != nil {
		log.Printf("Error al listar rutas: %v\n", err)
	}
}
