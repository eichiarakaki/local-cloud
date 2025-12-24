# Local Cloud System Design

## Overview
Local Cloud is a microservice-based system designed for downloading videos from URLs (primarily YouTube) within a Local Area Network (LAN). It allows users to submit video URLs through a web interface, download and process them server-side, store metadata in a MySQL database, and stream the videos via a React-based frontend. The system uses Go for backend services, TCP sockets for inter-service communication, and a custom message queue for orchestrating downloads asynchronously.

## Architecture Components

### 1. Web Server
- **Backend**: A Go server using Gorilla Mux for REST API endpoints. Runs on `0.0.0.0:3033`.
- **Frontend**: A React application built with Vite and Tailwind CSS. Served on port `3034`.
- **Responsibilities**:
  - Serve the web interface.
  - Handle user submissions of video URLs.
  - Query the database for video metadata.
  - Stream video files and thumbnails.

### 2. Message Queue
- A Go-based TCP broker running on `localhost:3032`.
- Acts as an intermediary for communication between the Web Server and Downloader Server.
- Manages a queue of download requests and tracks the status of the Downloader Server.
- Uses a custom FIFO linked list for queuing URLs.

### 3. Downloader Server
- A Go TCP server running on `localhost:3031`.
- Handles video downloads using `yt-dlp`, transcodes to MP4 with `ffmpeg`, generates thumbnails, and updates the database.
- Validates URLs (currently limited to YouTube).
- Communicates status via TCP protocol.

### 4. Database
- MySQL database for storing video metadata.
- Table: `localcloud` with fields for id, filepath, filename, thumbnail, created_at.
- Automatically created by the Downloader Server.

### 5. Shared Modules
- Go module for loading configuration from `config.json`.

## Communication Flow

### Workflow
1. **User Submission**:
   - User accesses the frontend and submits a video URL via the Downloader page.
   - Frontend sends a POST request to `/api/send-to-downloader-server` with the URL.

2. **Backend Processing**:
   - Web Server backend receives the URL and forwards it to the Message Queue via TCP connection.

3. **Queue Management**:
   - Message Queue enqueues the URL and checks the Downloader Server's status.
   - If the Downloader is free, it sends a `lock <URL>` command to initiate download.
   - Responds to backend with queue position and status.

4. **Download Process**:
   - Downloader Server validates the URL.
   - Downloads the video using `yt-dlp`.
   - Transcodes to MP4 and extracts thumbnail using `ffmpeg`.
   - Saves files to the configured storage path.
   - Inserts metadata into MySQL database.
   - Updates status to free for next request.

5. **Display**:
   - Frontend fetches video list from `/api/videos` and displays them.
   - Videos are streamed via `/api/videos-storage/{filename}`.

### Inter-Service Communication
- **Web Server ↔ Message Queue**: TCP connection for sending URLs and receiving status updates.
- **Message Queue ↔ Downloader Server**: TCP connection for status checks and issuing download commands.
- **Web Server ↔ Database**: Direct MySQL connections for querying metadata.
- **Downloader Server ↔ Database**: Direct MySQL connections for inserting metadata.
- **Frontend ↔ Backend**: HTTP requests for API calls and media streaming.

## Technologies Used
- **Backend Services**: Go (Golang)
- **Frontend**: React, Vite, Tailwind CSS
- **Database**: MySQL
- **Video Processing**: yt-dlp, ffmpeg
- **Communication**: TCP sockets, HTTP
- **Configuration**: JSON file

## Key Features
- Asynchronous downloads via message queue.
- Local storage and streaming.
- LAN-only access for security.
- Automatic thumbnail generation.
- MP4 transcoding for compatibility.

## Security and Resilience
- CORS configured for frontend origin.
- URL validation for supported domains.
- Error handling for DB failures (removes files on metadata insert failure).
- No authentication (LAN trust model).