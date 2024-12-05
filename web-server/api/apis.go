package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)


func GetVideoByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	videoID := vars["videoID"]
	testResponse := fmt.Sprintf("TEST RESPONSE %s", videoID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(testResponse))
}
