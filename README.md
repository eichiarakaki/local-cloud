# Local Cloud

**Local Cloud** is a microservice-based system designed to allow users within the same **Local Area Network (LAN)** to download videos via shared URLs and watch them through a locally hosted web application. 

> The system efficiently handles video downloads, stores the data, and renders them on a user-friendly front-end.


## Quick Start

### Using Docker Compose (Recommended)

1. **Clone and configure**:
   ```bash
   git clone <repo>
   cd LocalCloud
   ```

2. **Ensure `config.json` is set up** with Docker service names as hostnames (see `config.json` example below).

3. **Start all services**:
   ```bash
   docker-compose -f docker-compose-dev.yml up --build
   ```

4. **Access the application**:
   - Frontend: `http://localhost:3034`
   - Backend API: `http://localhost:3033/api`
   - MySQL will be ready after initial setup

**Example `config.json` for Docker**:
```json
{
  "video-storage-path": "/videos/",
  "mysql-conn": "admin:1010@tcp(mysql:3306)",
  "mysql-db-name": "localcloud",
  "mysql-table-name": "localcloud",
  "downloader-socket": "downloader-server:3031",
  "message-queue-socket": "message-queue:3032",
  "webserver-backend-socket": "web-backend:3033",
  "webserver-frontend-port": "3034"
}
```

### Manual Setup (Local Development)

If you prefer not to use Docker, you can run services locally:

**Dependencies**:
- Go 1.26.0+
- Node.js + npm or bun
- MySQL 9.6.0+
- `yt-dlp`: `pip install yt-dlp`
- `ffmpeg`: Install via your package manager

**Alternatively, use Nix**:
```bash
nix flake develop
```

**Then run each service** (in separate terminals):
```bash
cd downloader-server && go run .
cd message-queue && go run .
cd web-app/backend && go run .
cd web-app/frontend && npm install && npm run dev
```

**Configure frontend**: Ensure `.env` in `web-app/frontend/` matches your backend URL:
```
VITE_API_BASE=http://localhost:3033
```
## Previews
Home Page
![image](https://github.com/user-attachments/assets/fb5f8b1b-9f6c-4259-a6fd-5f55178f1573)
Download Page
![image](https://github.com/user-attachments/assets/b64a45d0-f6d0-49ff-bae3-edad9ca8e7af)


## **Key Features**

### **Web Server**  
- **Content Display**: Shows a list of videos stored in the database with the following metadata:  
  - File name  
  - File path (MP4)  
  - Automatically generated thumbnail
- **URL Submission**: Users can submit video URLs through an input field on the interface.

### **Downloader Server (Microservice)**
- **Download Management**: Handles video download requests, saves the video files, and generates related metadata.  
- **Database Updates**: After completing a download:
  - Saves the video file to the dedicated folder.
  - Generates and stores a thumbnail.
  - Updates the database with file path and metadata.  
- **TCP/IP Communication**: Maintains efficient communication with the Web Server and the Message Queue.

### **Message Queue**  
- **Inter-Service Communication**: Acts as an intermediary for asynchronous communication between the Web Server and the Downloader Server.  
- **Request Management**: Ensures that multiple download requests are processed efficiently and sequentially.

### **Database**  
- **Metadata Storage**: Stores relevant information about downloaded videos, including:  
  - File path (MP4)  
  - File name  
  - Thumbnail path  
- **Synchronization**: Ensures the front-end remains up-to-date with the downloaded videos.

### **Videos Folder**  
- **Local Storage**: All downloaded video files are saved in a designated directory accessible by the Downloader Server.


## **System Architecture**

Below is an overview of the **Local Cloud** system architecture:

![image](https://github.com/user-attachments/assets/07858769-96fc-4783-90f0-f4119d6d36d4)

### **Workflow**  
1. Users access the **Web Server** through a browser and submit a video URL.  
2. The Backend of the Web Server sends the URL to the **Message Queue** over a TCP/IP connection.  
3. The **Downloader Server** listens for messages from the queue and processes the video request.  
4. Once downloaded, the **Downloader Server** generates a thumbnail, stores the video file in the **Videos Folder**, and updates the **Database**.  
5. The **Web Server** queries the database to fetch video metadata and displays the available videos on the front-end.  


## **Technologies Used**

- **Backend**:  
  - Go (Golang) for the Web Server, Message Queue and Downloader Server.  
  - TCP/IP for inter-service communication.  
- **Frontend**:  
  - React + Vite + Tailwind CSS for the user interface.
- **Database**:  
  - MySQL for metadata storage.  
- **Containerization**:  
  - Docker and Docker Compose for development and deployment.  
- **Video Processing**:  
  - `yt-dlp` for video downloading.  
  - `ffmpeg` for MP4 transcoding and thumbnail extraction.

## **Technologies Used (Previous Section Content)**

- **Frontend**:  
  - React (Vite) and Tailwind for building the user interface.  
- **Database**:  
  - MySQL (configurable).  
- **File Storage**:  
  - Local file system for a simple folder.  



