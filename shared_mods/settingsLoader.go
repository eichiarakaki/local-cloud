package shared_mods

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var (
	VideoStoragePath       string
	MySQLConn              string
	MySQLDBName            string
	MySQLTableName         string
	DownloaderServerSocket string
	MessageQueueSocket     string
	WebServerSocket        string
)

// Must match with the JSON properties.
type Config struct {
	VideoStoragePath       string `json:"video-storage-path"`
	MySQLConn              string `json:"mysql-conn"`
	MySQLDBName            string `json:"mysql-db-name"`
	MySQLTableName         string `json:"mysql-table-name"`
	DownloaderServerSocket string `json:"downloader-socket"`
	MessageQueueSocket     string `json:"message-queue-socket"`
	WebServerSocket        string `json:"webserver-socket"`
}

// This function should be called only once in each main function of each library that wants to use these global variables.
func LoadConfig() error {
	// config.json MUST BE ON THE ROOT DIR.
	// This file should be on src for this to work propertly.
	file, err := os.Open("../config.json")
	if err != nil {
		return fmt.Errorf("Error while opening the config file: %s", err)
	}
	defer file.Close()

	// Reads the file content
	bytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Error when reading the config file: %s", err)
	}

	// Deserializing the JSON into a 'Config struct'
	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return fmt.Errorf("Error when deserializing the JSON file: %s", err)
	}

	VideoStoragePath = config.VideoStoragePath
	MySQLConn = config.MySQLConn
	MySQLDBName = config.MySQLDBName
	MySQLTableName = config.MySQLTableName
	DownloaderServerSocket = config.DownloaderServerSocket
	MessageQueueSocket = config.MessageQueueSocket
	WebServerSocket = config.WebServerSocket

	return nil
}
