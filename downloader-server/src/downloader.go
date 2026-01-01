package src

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	shared "shared_mods"
	"strings"
)

func StartDownload(url string, conn net.Conn) {
	_, err := conn.Write([]byte(ServerStatus))
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("INFO: Processing %s\n", url)

	Download(url)
	mu.Lock()
	ServerStatus = Free // Changes Status to free server
	mu.Unlock()
	_, err = conn.Write([]byte(ServerStatus))
	if err != nil {
		log.Println("[ERROR]", err)
		return
	} // Sends the status to the frontend immediately after the download ends
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
	Title     string
	Path      string
	Thumbnail string
}

func Download(url string) {
	cmd := commandBuilder(url)

	// Executes the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error when executing yt-dlp: %v\n", err)
	}

	videoData, err := dbWrapper(output)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	// log.Printf("%s\n", output)
	log.Printf("Download Completed: %s\nFile name: %s\n", videoData.Path, videoData.Title)

	videoData.PushToBD()
}

func commandBuilder(url string) *exec.Cmd {
	// To-Do: give the user options to download playlists, and override files
	outputTemplate := shared.VideoStoragePath + "%(title)s.%(ext)s"
	return exec.Command(
		"yt-dlp",
		url,
		"--no-playlist",
		"--embed-thumbnail",
		"-f", "bestvideo[height<=1080]+bestaudio/best[height<=1080]",
		"--output", outputTemplate,
	)
}

func dbWrapper(output []byte) (*VideoData, error) {
	var data VideoData

	fullfilepath, filename, isDownloaded, err := extractTitleAndPath(output)
	if err != nil {
		log.Fatalln(err)
	}

	thumbnailPath, err := extractThumbnail(fullfilepath, filename)
	if err != nil {
		log.Fatalln(err)
	}

	// Converting from any format to mp4
	fullfilepath, err = MP4Transformer(fullfilepath)
	if err != nil {
		log.Fatalf("WARNING: Couldn't convert %s to mp4: %s\n", fullfilepath, err)
	}

	data.Path = fullfilepath
	data.Title = filename
	data.Thumbnail = thumbnailPath

	if isDownloaded {
		fmt.Println("Video already exists in the database.")
	}

	return &data, nil
}

// Returns Full File Path (including the filename), File Name (only), if the file already exists, and an Error.
func extractTitleAndPath(output []byte) (string, string, bool, error) {
	// Filters the OUTPUT
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if strings.Contains(line, "has already been downloaded") {
			parts := strings.Split(line, "has already been downloaded")
			ft := strings.TrimSpace(parts[0])
			fullPath := ft[11:] // removes the "[download]" part from the output
			fileName := filterVideoTitle(fullPath)

			return fullPath, fileName, true, nil
		}
		if strings.Contains(line, "[Merger] Merging formats into ") {
			parts := strings.Split(line, "formats into ") // Splits when finds a `[Merger] Merging formats into "`
			if len(parts) == 2 {
				// Filtering the full path
				fullPath := strings.TrimSpace(parts[1])
				fullPath = fullPath[1 : len(fullPath)-1] // Removes the "" from the output
				fileName := filterVideoTitle(fullPath)

				return fullPath, fileName, false, nil
			}
		}
	}
	return "", "", false, fmt.Errorf("Error extracting title and path.")
}

func filterVideoTitle(fullPath string) string {
	// Filtering the full path from the title
	fileTypes := []string{".mp4", ".mp3", ".webm", ".avi", ".flv", ".mkv"}
	fileName := strings.Split(fullPath, shared.VideoStoragePath)[1]

	for _, fileType := range fileTypes { // Removes the file extension from the title.
		if strings.HasSuffix(fileName, fileType) {
			fileName = strings.Split(fileName, fileType)[0]
		}
	}
	return fileName
}

// MP4Transformer Transforms to MP4 format to have compatibility with IOS devices, since IOS does not have support for mkv in web browsers like Safari.
func MP4Transformer(inputPath string) (string, error) {
	log.Println("INFO: Transforming file to MP4 format...")
	// Extract the folder and the file name
	dir := filepath.Dir(inputPath)
	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath)) // Gets the filename without the extension.

	// Create the output file path with the .mp4 extension.
	outputPath := filepath.Join(dir, baseName+".mp4")

	// FFmpeg command for converting the file.
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputPath, // Input file

		"-vcodec", "libx264", // Codification video en H.264

		/*
		 Options:
		 ultrafast: Fast but worse file compression.
		 Faster: Faster than average, but without sacrificing too much quality.
		 Medium: Balance between speed and quality (default in FFmpeg).
		 Slow: Slow, but better file compression.
		 Veryslow: Very slow, but optimal compression.
		*/
		"-preset", "faster",

		// Range: 0-51. Recommended range: 18-23.
		"-crf", "18", // Constant rate factor for similar quality

		"-acodec", "aac", // AAC audio encoding

		// 128k (acceptable audio quality).
		// 192k (high audio quality).
		// 320k (max audio quality).
		"-b:a", "192k", // Audio bitrate similar to the original

		"-movflags", "+faststart", // Move moov atom to beginning for web streaming

		outputPath,
	)

	// Executes the command
	err := cmd.Run()
	if err != nil {
		return inputPath, err // If there's an error, just returns que input file.
	}

	// Removing the input file (already downloaded) since it will no longer be used.
	err = os.Remove(inputPath)
	if err != nil {
		log.Printf("WARNING: %s removed unsuccessfully. Remove it by yourself.", err)
	} else {
		log.Printf("INFO: %s removed successfully.", inputPath)
	}

	return outputPath, nil
}

// Creates the thumbnail if it doesn't exist, and only returns the full path to the thumbnail if exists.
func extractThumbnail(fullpath string, filename string) (string, error) {
	thumbnailTmpl := fmt.Sprintf("%s%s_thumbnail.webp", shared.VideoStoragePath, filename)
	if _, err := os.Stat(thumbnailTmpl); os.IsNotExist(err) { // If thumbnailTmpl does NOT exists.
		cmd := exec.Command("ffmpeg", "-dump_attachment:t:0", thumbnailTmpl, "-i", fullpath, "-vn", "-f", "null", "-")
		// Executes the command
		output, err := cmd.CombinedOutput()
		fmt.Println(string(output))
		if err != nil {
			log.Fatalf("Error when extracting thumbnail: %v\n", err)
		}
		log.Println("Thumbnail was extracted successfully.")
	} else {
		log.Println("Thumbnail already exists.")
	}
	return thumbnailTmpl, nil
}
