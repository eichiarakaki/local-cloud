package main

import (
    "fmt"
    "log"
    "os/exec"
)

func main() {
    url := "https://www.youtube.com/watch?v=9862-aHQ8zY"
    outputFileName := "video.mp4"

    cmd := exec.Command("yt-dlp", "-o", outputFileName, url)

    // Executes the command
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatalf("Error when executing yt-dlp: %v\n", err)
    }

    fmt.Printf("yt-dlp output: %s\n", output)
    fmt.Println("Download Completed")
}

