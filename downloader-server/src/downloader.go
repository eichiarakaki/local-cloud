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

/*
--live-from-start Download livestreams from the start.  Currently only supported for YouTube (Experimental)
--socket-timeout SECONDS Time to wait before giving up, in seconds
--min-filesize SIZE Abort download if filesize is smaller than SIZE, e.g.  50k or 44.6M
--max-filesize SIZE Abort download if filesize is larger than SIZE, e.g.  50k or 44.6M
--yes-playlist Download the playlist, if the URL refers to a video and a playlist
--force-overwrites Overwrite all video and metadata files.  This option includes --no-continue
--write-info-json Write video metadata to a .info.json file (this may contain personal information)
--load-info-json FILE JSON file containing the video information (created with the "--write-info-json" option)
--embed-thumbnail Embed thumbnail in the video as cover art
-q, --quiet Activate quiet mode.  If used with --verbose, print the log to stderr
--progress Show progress bar, even if in quiet mode
--write-subs Write subtitle file
*/

func Download(url string) {

	cmd := exec.Command("yt-dlp", url, "--yes-playlist", "--force-overwrites", "--quiet", "--progress")

	// Executes the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error when executing yt-dlp: %v\n", err)
	}

	fmt.Printf("yt-dlp output: %s\n", output)
	fmt.Println("Download Completed")
}
