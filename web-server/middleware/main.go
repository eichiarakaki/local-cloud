package middleware

import (
	"log"
	"net/http"
)

// Allowed paths
var allowedPaths = map[string]bool{
	"/":     true,
	"/api/": false,
}

func APIFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		if allowedPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
		} else if allowedPaths[r.URL.Path] == false {
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
		} else {
			http.Error(w, "URL Unavailable.", http.StatusServiceUnavailable)
		}
	})
}
