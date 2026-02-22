# Local Cloud Interview Questions and Answers

## 1. Can you describe the overall architecture of your Local Cloud project?
**Answer**: Local Cloud is a microservice-based system consisting of four main components: a Web Server (with backend and frontend), a Message Queue, a Downloader Server, and a MySQL database. The backend and other services are built in Go, while the frontend uses ReactJS with Vite and Tailwind CSS. The system allows users on a LAN to submit video URLs, download them server-side, and stream them via a web interface. Services communicate via TCP sockets, with the Message Queue acting as an orchestrator for asynchronous downloads.

## 2. How do the different services in your system communicate with each other?
**Answer**: The services use TCP sockets for inter-service communication. The Web Server backend connects to the Message Queue to submit download requests. The Message Queue maintains a TCP connection to the Downloader Server to check its status and issue download commands. The Web Server and Downloader Server both connect directly to the MySQL database for reading and writing metadata. The frontend communicates with the backend via HTTP REST API calls.

## 3. Why did you choose Go for building the backend services?
**Answer**: Go was chosen for its excellent performance in concurrent operations, which is crucial for handling multiple download requests and TCP connections. It's also well-suited for microservices due to its simplicity, strong standard library for networking, and efficient compilation. Additionally, Go's goroutines make it easy to handle asynchronous tasks like video processing without complex threading models.

## 4. What is the role of the Message Queue in your application?
**Answer**: The Message Queue serves as an intermediary broker between the Web Server and the Downloader Server. It maintains a queue of video URLs submitted by users and ensures downloads are processed sequentially. It tracks the Downloader Server's status (free/busy) and only issues download commands when the server is available, preventing conflicts and providing feedback on queue positions to users.

## 5. Can you walk me through the workflow from when a user submits a video URL to when the video becomes available for watching?
**Answer**: First, the user enters a URL on the frontend's Downloader page, which sends a POST request to the backend API. The backend forwards the URL to the Message Queue via TCP. The Message Queue enqueues the request and, when the Downloader Server is free, sends a 'lock <URL>' command. The Downloader validates the URL, downloads the video using yt-dlp, transcodes it to MP4 with ffmpeg, generates a thumbnail, saves files to storage, and inserts metadata into the database. Finally, the frontend can fetch and display the new video from the updated database.

## 6. How does your system handle video processing and storage?
**Answer**: Video processing is handled by the Downloader Server using yt-dlp for downloading and ffmpeg for transcoding to MP4 format and thumbnail extraction. Videos and thumbnails are stored locally in a configured directory. The system ensures browser compatibility by converting to MP4 and generates thumbnails automatically. If database insertion fails, the files are cleaned up to maintain consistency.

## 7. What challenges did you face in building this microservice system, and how did you overcome them?
**Answer**: One challenge was coordinating asynchronous communication between services over TCP sockets, which we solved by implementing a custom protocol with status messages. Another was ensuring data consistency between file storage and database metadata, addressed by defensive cleanup on errors. Handling concurrent downloads was managed through the Message Queue's status tracking. We also faced issues with video format compatibility, resolved by mandatory MP4 transcoding.

## 8. How is the database integrated into your system, and what kind of data do you store?
**Answer**: The database is MySQL, integrated directly into the Web Server backend for queries and the Downloader Server for inserts. We store video metadata including file paths, filenames, thumbnail paths, and creation timestamps. The table is automatically created by the Downloader Server on startup. The backend queries this data to populate the frontend, while the Downloader updates it after successful downloads.

## 9. Describe how the frontend interacts with the backend.
**Answer**: The frontend is a React SPA that makes HTTP requests to the backend's REST API. It fetches video lists via GET /api/videos, retrieves individual video data via GET /api/video/{id}, and streams media files via GET /api/videos-storage/{filename}. For downloads, it posts URLs to POST /api/send-to-downloader-server. The backend serves the frontend statically and handles CORS to allow the configured origin.

## 10. How do you ensure the system is secure, especially since it's designed for LAN use?
**Answer**: Security is based on the LAN trust model, with no authentication required. We implement CORS to restrict API access to the configured frontend origin. URL validation limits downloads to specific domains (currently YouTube). The system doesn't expose services externally, relying on local network isolation. Database connections use configured credentials, and file storage is local with proper path handling to prevent directory traversal.