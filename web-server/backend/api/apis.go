package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	shared "shared_mods"
	"strings"

	"github.com/eichiarakaki/local-cloud/web-server/backend/src"
	"github.com/gorilla/mux"
)

type VideoData struct {
	ID        int    `json:"id"`
	Path      string `json:"filepath"`
	Title     string `json:"filename"`
	Thumbnail string `json:"thumbnail"`
	CreatedAt string `json:"created_at"`
}

// Processes and returns the relative file paths, thumbnails and titles stored in videos-storage.
func GetAllVideos(w http.ResponseWriter, r *http.Request) {
	db, err := src.ConnectDB()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	cmd := fmt.Sprintf("SELECT id, filepath, filename, thumbnail, created_at FROM %s", shared.MySQLTableName)
	rows, err := db.Query(cmd)
	if err != nil {
		// In most cases, this error occurs because the table does not exist in the database, either way it returns an empty json to the frontend.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("[]"))
		if err != nil {
			log.Println("ERROR: writing empty JSON response:", err)
		}
		return
	}
	defer rows.Close()

	// This handles when there's a table but it is an EMPTY table.
	if !rows.Next() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("[]"))
		if err != nil {
			log.Println("ERROR: writing empty JSON response:", err)
		}
		return
	}

	// Process the query results
	var videosData []VideoData
	for {
		var aVideoData VideoData
		if err := rows.Scan(
			&aVideoData.ID,
			&aVideoData.Path,
			&aVideoData.Title,
			&aVideoData.Thumbnail,
			&aVideoData.CreatedAt); err != nil {
			log.Println(err)
			http.Error(w, "Error processing the data", http.StatusInternalServerError)
			return
		}
		// Decoding Path and Thumbnail to handle spaces and special characters
		decodedThumbnail, err := url.QueryUnescape(filepath.Base(aVideoData.Thumbnail))
		if err != nil {
			log.Println("Error decoding thumbnail URL:", err)
		}
		decodedPath, err := url.QueryUnescape(filepath.Base(aVideoData.Path))
		if err != nil {
			log.Println("Error decoding video file URL:", err)
		}

		aVideoData.Path = decodedPath
		aVideoData.Thumbnail = decodedThumbnail

		videosData = append(videosData, aVideoData)

		if !rows.Next() { // Exits the loop if there a re no more rows
			break
		}
	}

	log.Println("INFO: Got a request '/api/videos'.")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(videosData); err != nil {
		log.Println(err)
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

// Returns a JSON for rendering in the frontend.
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

// Renders the images and videos (by encoded name).
func ServeStorage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mediaName := vars["mediaName"]

	// Media's full path
	path := shared.VideoStoragePath + mediaName

	// Opening file
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, "Media not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Getting media's info
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Error reading file info!", http.StatusInternalServerError)
		return
	}

	// Determine the content type
	ext := strings.ToLower(filepath.Ext(mediaName))
	var contentType string
	switch ext {
	case ".mp4":
		contentType = "video/mp4"
	case ".webm":
		contentType = "video/webm"
	case ".mkv":
		contentType = "video/x-matroska"
	case ".webp":
		contentType = "image/webp"
	default:
		contentType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", contentType)
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

func Testing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["id"]

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(ID)
	if err != nil {
		http.Error(w, "Error enconding response", http.StatusInternalServerError)
	}
}