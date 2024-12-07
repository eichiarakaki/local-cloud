package src

/*
To-Do
- Handle unexpected errors and send to the frontend
*/
import (
	"fmt"
	"log"
	"os/exec"
	shared "shared_mods"
	"strings"
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
-S, --format-sort SORTORDER Sort the formats by the fields given, see "Sorting Formats" for more details
*/

type VideoData struct {
	Title string
	Path  string
}

func Download(url string) {
	cmd := commadBuilder(url)

	// Executes the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error when executing yt-dlp: %v\n", err)
	}

	videoData, isDownloaded, err := dbWrapper(output)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	// log.Printf("%s\n", output)
	log.Printf("Download Completed: %s\nFile name: %s\n", videoData.Path, videoData.Title)

	if !isDownloaded {
		videoData.PushToBD()
	} else {
		// To-Do: notify the frontend that the file already exists in the database.
		log.Println("File already exists in the database.")
	}
}

func commadBuilder(url string) *exec.Cmd {
	// To-Do: give the user options to download playlists, and override files
	outputTemplate := shared.VideoStoragePath + "%(title)s.%(ext)s"
	return exec.Command("yt-dlp", url, "--no-playlist", "--output", outputTemplate)
}

func dbWrapper(output []byte) (*VideoData, bool, error) {
	var data VideoData

	// Filters the OUTPUT
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "has already been downloaded") {
			return &data, true, nil
		}
		if strings.Contains(line, "[Merger] Merging formats into ") {
			parts := strings.Split(line, "formats into ") // Splits when finds a `[Merger] Merging formats into "`
			if len(parts) == 2 {
				// Filtering the full path
				fullPath := strings.TrimSpace(parts[1])
				fullPath = fullPath[1 : len(fullPath)-1] // Removes the "" from the output

				// Filtering the full path from the title
				fileTypes := []string{".mp4", ".mp3", ".webm", ".avi", ".flv", ".mkv"}
				fileName := strings.Split(fullPath, shared.VideoStoragePath)[1]

				for _, fileType := range fileTypes { // Removes the file extension from the title.
					if strings.HasSuffix(fileName, fileType) {
						fileName = strings.Split(fileName, fileType)[0]
					}
				}

				data.Path = fullPath
				data.Title = fileName
				return &data, false, nil // removes the "" and returns it
			}
		}
	}
	return nil, false, fmt.Errorf("Error on database's wrapper.")
}
