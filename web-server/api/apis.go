package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	shared "shared_mods"

	"github.com/eichiarakaki/local-cloud/src"
	"github.com/gorilla/mux"
)

type VideoData struct {
	ID        int    `json:"id"`
	Path      string `json:"filepath"`
	Title     string `json:"filename"`
	Thumbnail string `json:"thumbnial"`
	CreatedAt string `json:"created_at"`
}

func GetAllVideos(w http.ResponseWriter, r *http.Request) {
	db, err := src.ConnectDB()
	if err != nil {
		log.Println(err)
	}

	cmd := fmt.Sprintf("SELECT id, filepath, filename, thumbnail, created_at FROM %s", shared.MySQLTableName)
	rows, err := db.Query(cmd)
	if err != nil {
		http.Error(w, "Error when consulting the database.", http.StatusInternalServerError)
	}
	defer rows.Close()

	// Process
	var videosData []VideoData
	for rows.Next() {
		var aVideoData VideoData
		if err := rows.Scan(&aVideoData.ID, &aVideoData.Path, &aVideoData.Title, &aVideoData.Thumbnail, &aVideoData.CreatedAt); err != nil {
			http.Error(w, "Error processing the data", http.StatusInternalServerError)
			return
		}
		videosData = append(videosData, aVideoData)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videosData)
}

func GetVideoByID(w http.ResponseWriter, r *http.Request) {
	db, err := src.ConnectDB()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	vars := mux.Vars(r)
	videoID := vars["videoID"]
	var videodata VideoData

	query := fmt.Sprintf("SELECT id, filepath, filename, thumbnail, created_at FROM %s WHERE id = ?", shared.MySQLTableName)
	err = db.QueryRow(query, videoID).Scan(&videodata.ID, &videodata.Path, &videodata.Title, &videodata.Thumbnail, &videodata.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Video Not Found.", http.StatusNotFound)
		} else {
			http.Error(w, "Error when consulting.", http.StatusInsufficientStorage)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(videodata)
	if err != nil {
		http.Error(w, "Error enconding response", http.StatusInternalServerError)
	}
}
