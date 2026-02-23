package api

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
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
	"time"

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

// GetAllVideos Processes and returns the relative file paths, thumbnails and titles stored in videos-storage.
func GetAllVideos(w http.ResponseWriter, r *http.Request) {
	db, err := src.ConnectDB()
	if err != nil {
		log.Println(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("[ERROR] Closing db connection")
		}
	}(db)

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
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("[ERROR] Closing rows:", err)
		}
	}(rows)

	// This handles when there's a table, but it is an EMPTY table.
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

		if !rows.Next() { // Exits the loop if there are no more rows
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

// GetVideoByID Returns a JSON for rendering in the frontend.
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
		if errors.Is(err, sql.ErrNoRows) {
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

	// URL decode the mediaName
	mediaName, err := url.QueryUnescape(mediaName)
	if err != nil {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3034")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

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
	log.Println("[INFO] Got a request to '/api/SendToDownloaderServer/'")
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

	// Parsing the JSON received
	var requestBody struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Validates the URL
	if requestBody.URL == "" {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Sending the data to the message queue
	res, err := sendToMQ(shared.MessageQueueSocket, requestBody.URL)
	if err != nil {
		log.Println("[ERROR] Sending message to server:", err, "\n")
		http.Error(w, "Failed to send message to server:", http.StatusInternalServerError)
		return
	}
	fmt.Println("[BACKEND] Got", res)

	// Response to the frontend
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func Testing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["id"]

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(ID)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func sendToMQ(address string, message string) (string, error) {
	var conn net.Conn
	var err error
	for i := 0; i < 5; i++ { // Tries to connect to the socket for 5 times before giving up, in case the message queue server is not ready yet.
		conn, err = net.Dial("tcp", address)
		if err == nil {
			fmt.Printf("[INFO] Connected to message queue at %s\n", address)
			break
		}
		log.Printf("[ERROR] Failed to connect to socket (attempt %d/5): %v\n", i+1, err)
		time.Sleep(5 * time.Second) // Waits for 5 seconds before retrying
	}
	if conn == nil {
		return "", fmt.Errorf("[ERROR] Failed to connect to socket after 5 attempts")
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("[ERROR] Failed to close socket:", err)
		}
	}(conn)

	_, err = conn.Write([]byte(message))
	if err != nil {
		return "", fmt.Errorf("[ERROR] Failed to send to socket: %w", err)
	}

	// Reads response from the message queue for sending to the frontend.
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return "", nil
		}
		return "", fmt.Errorf("error when reading from the server: %s", err)
	}

	return response, nil
}
