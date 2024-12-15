package shared_mods

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	VideoStoragePath       string
	MySQLConn              string
	MySQLDBName            string
	MySQLTableName         string
	DownloaderServerSocket string
	MessageQueueSocket     string
	WebServerBackendSocket string
	WebServerFrontendPort  string
)

// Config Must match with the JSON properties.
type Config struct {
	VideoStoragePath       string `json:"video-storage-path"`
	MySQLConn              string `json:"mysql-conn"`
	MySQLDBName            string `json:"mysql-db-name"`
	MySQLTableName         string `json:"mysql-table-name"`
	DownloaderServerSocket string `json:"downloader-socket"`
	MessageQueueSocket     string `json:"message-queue-socket"`
	WebServerBackendSocket string `json:"webserver-backend-socket"`
	WebServerFrontendPort  string `json:"webserver-frontend-port"`
}

// LoadConfig This function should be called only once in each main function of each library that wants to use these global variables.
func LoadConfig(configFilePath string) error {
	// config.json MUST BE ON THE ROOT DIR.
	file, err := os.Open(configFilePath)
	if err != nil {
		return fmt.Errorf("error while opening the config file: %s", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	// Reads the file content
	bytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error when reading the config file: %s", err)
	}

	// Deserializing the JSON into a 'Config struct'
	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return fmt.Errorf("error when deserializing the JSON file: %s", err)
	}

	VideoStoragePath = config.VideoStoragePath
	MySQLConn = config.MySQLConn
	MySQLDBName = config.MySQLDBName
	MySQLTableName = config.MySQLTableName
	DownloaderServerSocket = config.DownloaderServerSocket
	MessageQueueSocket = config.MessageQueueSocket
	WebServerBackendSocket = config.WebServerBackendSocket
	WebServerFrontendPort = config.WebServerFrontendPort

	return nil
}
