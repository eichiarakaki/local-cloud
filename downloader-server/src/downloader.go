package src

import (
	"fmt"
	"log"
	"os/exec"
)

func StartDownload(url string) {
	mu.Lock()
	ServerStatus = Busy // Changes status to busy server
	mu.Unlock()

	go func() {
		Download(url)
		mu.Lock()
		ServerStatus = Free // Changes Status to free server
		mu.Unlock()
	}()
}

func Download(url string) {
	outputFileName := "new.mp4"

	cmd := exec.Command("yt-dlp", "-o", outputFileName, url)

	// Executes the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error when executing yt-dlp: %v\n", err)
	}

	fmt.Printf("yt-dlp output: %s\n", output)
	fmt.Println("Download Completed")
}
