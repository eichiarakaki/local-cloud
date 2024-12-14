package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
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

	log.Println("INFO: Got a request to '/api/videos'.")

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

	log.Printf("INFO: Got a request to '/api/video/%s'.", videoID)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(videodata)
	if err != nil {
		http.Error(w, "Error enconding response", http.StatusInternalServerError)
	}
}

// ServeStorage Renders the images and videos (by encoded name).
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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("[ERROR] Writing media file to disk:", err)
		}
	}(file)

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

// SendToDownloaderServer Handles the request from the frontend and sends to the message queue socket.
func SendToDownloaderServer(w http.ResponseWriter, r *http.Request) {
	// Configuring the CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3034")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Managing preflight request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Only allows POST Method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Reading the requested body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[ERROR] Reading the body request:", err)
		http.Error(w, "Failed to read request body:", http.StatusInternalServerError)
		return
	}

	// Parsing the JSON received
	type InputData struct {
		URL string `json:"url"`
	}

	var input InputData
	if err := json.Unmarshal(body, &input); err != nil {
		log.Println("[ERROR] Parsing the body request:", err)
		http.Error(w, "Failed to parse the body request:", http.StatusInternalServerError)
		return
	}

	// Validates the URL
	if input.URL == "" {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Sending the data to the message queue
	if err := sendToSocket(shared.MessageQueueSocket, input.URL); err != nil {
		log.Println("[ERROR] Sending message to server:", err)
		http.Error(w, "Failed to send message to server:", http.StatusInternalServerError)
		return
	}
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

func sendToSocket(address string, message string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to connect to socket: %w", err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("[ERROR] Failed to close socket:", err)
		}
	}(conn)

	_, err = conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to send to socket: %w", err)
	}

	return nil
}
