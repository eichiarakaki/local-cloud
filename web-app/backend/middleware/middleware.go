package middleware

import (
	"fmt"
	"log"
	"net/http"
	shared "shared_mods"
	"strings"
)

// Allowed paths
var allowedPaths = map[string]bool{
	"/":     true,
	"/api/": false,
}

func APIFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		log.Println(r.UserAgent(), r.RemoteAddr)

		// Check if the request is from the server
		if strings.HasPrefix(r.URL.Path, "/static/") {
			referer := r.Referer()
			socketURL := fmt.Sprintf("http://%s/", shared.WebServerBackendSocket)

			if referer == "" || !strings.HasPrefix(referer, socketURL) {
				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		// Handle the allowed paths for users
		if allowedPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
		} else if allowedPaths[r.URL.Path] == false {
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
		} else {
			http.Error(w, "URL Unavailable.", http.StatusServiceUnavailable)
		}
	})
}
