# Local Cloud — System Design and Developer Skills

## Executive Summary
Local Cloud is a LAN-first microservice system for submitting video URLs, downloading and transcoding them server-side, storing metadata in MySQL, and presenting the catalog via a React web app. Services communicate over TCP sockets with a lightweight custom message queue orchestrating downloads. The stack centers on Go for backend services and React + Vite + Tailwind for the frontend.

## Architecture Overview
- Services
  - Web Server Backend (`Go`, REST on `gorilla/mux`) at `web-server/backend/main.go:35`.
  - Frontend (`React` + `Vite` + `Tailwind`) at `web-server/frontend/` served on `config.json:webserver-frontend-port`.
  - Message Queue (`Go`, TCP broker) at `message-queue/main.go:88-113`.
  - Downloader Server (`Go`, TCP server) at `downloader-server/src/server.go:26-53`.
  - Shared config loader (`Go`) at `shared_mods/settingsLoader.go:34-69` backed by `config.json`.
  - MySQL database seeded by `downloader-server/sql/createTable.sql:1-9`.
- Data flow
  - User submits a URL in the frontend (`web-server/frontend/src/pages/Downloader.tsx:101-116`).
  - Backend receives and forwards URL to the Message Queue (`web-server/backend/api/apis.go:197-249`, `apis.go:263-290`).
  - Message Queue tracks status and enqueues URLs (`message-queue/queue.go:29-47`, `queue.go:49-57`).
  - When downloader is free, MQ issues a `lock <URL>` to the Downloader Server (`message-queue/main.go:66-76`).
  - Downloader validates URL and downloads via `yt-dlp`, then transcodes to MP4 and extracts thumbnail with `ffmpeg` (`downloader-server/src/downloader.go:75-86`, `downloader.go:160-216`, `downloader.go:219-234`).
  - Downloader stores metadata in MySQL (`downloader-server/src/database.go:100-127`).
  - Frontend lists videos from backend `/api/videos` and streams via `/api/videos-storage/:name` (`web-server/frontend/src/pages/Home.tsx:18-27`, `web-server/backend/api/router.go:10-18`, `web-server/backend/api/apis.go:149-195`).

## Components
- Web Server Backend
  - Bootstraps `mux` router and CORS (`web-server/backend/main.go:22-43`).
  - Loads config (`shared_mods.LoadConfig`) early (`web-server/backend/main.go:15-21`).
  - API endpoints (`web-server/backend/api/router.go:7-19`):
    - `GET /api/videos` → list metadata (`web-server/backend/api/apis.go:31-115`).
    - `GET /api/video/{videoID}` → single record (`apis.go:117-147`).
    - `GET /api/videos-storage/{mediaName}` → static streaming (`apis.go:149-195`).
    - `POST /api/send-to-downloader-server` → forwards URL to MQ (`apis.go:197-249`).
  - DB connectivity (`web-server/backend/src/database.go:13-29`) via DSN from `config.json`.
  - Storage path joins via `shared.VideoStoragePath` (`web-server/backend/api/apis.go:155`).

- Frontend
  - Routes (`web-server/frontend/src/App.tsx:12-27`): home, single video, browser, downloader.
  - Home fetches `GET /api/videos` and renders `VideoCard` (`web-server/frontend/src/pages/Home.tsx:18-27`, `web-server/frontend/src/elements/videoElement.tsx:8-27`).
  - Player streams via backend storage endpoint (`web-server/frontend/src/pages/SingleVideoPage.tsx:23-27`).
  - Browser page filters client-side (`web-server/frontend/src/elements/Browser.tsx:35-41`).
  - Downloader posts URL JSON (`web-server/frontend/src/pages/Downloader.tsx:26-36`) and displays broker status.

- Message Queue
  - Maintains TCP connection to Downloader Server and tracks its status (`message-queue/main.go:28-60`).
  - Enqueues backend URLs and responds with JSON status (`message-queue/main.go:115-153`, `message-queue/utils/messenger.go:1-16`).
  - Dequeues when status becomes `free\n` (`message-queue/main.go:66-76`).
  - Custom FIFO linked list queue (`message-queue/queue.go:16-27`, `queue.go:29-47`, `queue.go:49-57`).

- Downloader Server
  - TCP server exposes status-first protocol and handles commands (`downloader-server/src/server.go:55-126`).
  - Validates URL prefixes (YouTube only) (`downloader-server/src/urlFilter.go:8-22`).
  - Orchestrates download, MP4 transform, thumbnail extraction, and DB upload (`downloader-server/src/downloader.go:14-31`, `downloader.go:55-74`, `downloader.go:160-216`, `downloader.go:219-234`, `downloader-server/src/database.go:100-127`).
  - Defensive cleanup on DB failures (`downloader-server/src/database.go:84-99`).

- Configuration and Secrets
  - Centralized in `config.json` (`config.json:1-10`) loaded by `shared_mods/settingsLoader.go:34-69`.
  - Key fields: `video-storage-path`, `mysql-conn`, `mysql-db-name`, `mysql-table-name`, socket addresses, and frontend port.

## API Surface
- Backend
  - `GET /api/videos` → `VideoData[]` with `{ id, filepath, filename, thumbnail, created_at }` (`web-server/backend/api/apis.go:23-29`, `apis.go:108-115`).
  - `GET /api/video/{videoID}` → `VideoData` (`web-server/backend/api/apis.go:117-147`).
  - `GET /api/videos-storage/{mediaName}` → binary content (`video/mp4`, `image/webp`, etc.) (`web-server/backend/api/apis.go:177-195`).
  - `POST /api/send-to-downloader-server` body `{ url }` → broker JSON `MQBackend` (`message-queue/utils/messenger.go:4-8`).

- Message Queue → Backend response (`MQBackend`)
  - `server_status`: downloader state (`"free\n"`, `"busy\n"`, `"inurl\n"`).
  - `queue_position`: 1-based index including active download when busy (`message-queue/queue.go:93-107`).
  - `message`: enqueue confirmation (`message-queue/main.go:132-148`).

- Downloader Server protocol (TCP)
  - Client reads initial status (`downloader-server/src/server.go:56-59`).
  - Commands: `test <URL>` for validation; `lock <URL>` to start a download (`server.go:86-124`).
  - Status messages: `free\n`, `busy\n`, `inurl\n` (`server.go:17-25`).

## Data Model
- Table `localcloud` (name configurable):
  - `id` (INT, PK, AUTO_INCREMENT)
  - `filepath` (VARCHAR(250))
  - `filename` (VARCHAR(100))
  - `thumbnail` (VARCHAR(350))
  - `created_at` (TIMESTAMP DEFAULT CURRENT_TIMESTAMP)
  - DDL at `downloader-server/sql/createTable.sql:1-9`.

## Storage and Media Handling
- Downloads saved under `config.json:video-storage-path`.
- Transcoded to `mp4` for browser/device compatibility (`downloader-server/src/downloader.go:160-216`).
- Thumbnail extraction from container metadata (`downloader-server/src/downloader.go:219-234`).
- Backend streams by filename via `/api/videos-storage/:mediaName` (`web-server/backend/api/apis.go:149-195`).

## Error Handling and Resilience
- DB connection validated with `Ping()` before use (`web-server/backend/src/database.go:22-28`, `downloader-server/src/database.go:24-31`).
- Empty table and missing table handled gracefully in `GET /api/videos` returning `[]` (`web-server/backend/api/apis.go:63-72`).
- Invalid URLs flagged and status cycles reset (`downloader-server/src/server.go:103-118`).
- On DB write errors, media files are removed to avoid drift (`downloader-server/src/database.go:84-99`).
- CORS configured to allow frontend origin (`web-server/backend/api/apis.go:201-204`).

## Security Considerations
- Local network trust; no authentication on API or TCP sockets.
- CORS restricts origin to configured frontend (currently `http://localhost:3034`) (`web-server/backend/api/apis.go:201`).
- URL validation limited to YouTube domains (`downloader-server/src/urlFilter.go:8-22`).

## Development Setup

### Docker Compose (Recommended)
The project uses `docker-compose-dev.yml` for full local development with containerized services:

**Prerequisites**:
- Docker and Docker Compose installed
- `config.json` configured in the project root with appropriate paths and socket addresses

**Services**:
- **web-frontend**: React + Vite dev server on port `3034`
- **web-backend**: Go REST API on port `3033` with live reload via `go run`
- **message-queue**: Go TCP broker on port `3032`
- **downloader-server**: Go TCP server on port `3031` with `yt-dlp` and `ffmpeg`
- **mysql**: MySQL 9.6.0 database with persistent volume

**Build Context**: Repository root (`.`); services build with mounted workspaces using `go.work` for local module resolution.

**Run**:
```bash
docker-compose -f docker-compose-dev.yml up --build
```

Data persists in the `mysql_data` volume and `./videos` directory.

### Local Development (Nix / Manual)
Alternatively, use the provided `flake.nix` for a Nix dev shell or manually install:
- Go 1.26.0+
- Node.js + npm/bun
- MySQL 9.6.0+
- `yt-dlp` and `ffmpeg`

Then run services individually from their respective directories using `go run .` or `npm run dev`.

## Deployment and Operations
- Dependencies: MySQL, `yt-dlp`, `ffmpeg`, Node.js (`README.md:11-18`).
- Startup order: Downloader → Message Queue → Web Server Backend → Frontend (`README.md:18`).
- Sockets and ports from `config.json` (`config.json:6-9`).
- Observability: structured logs across services; route printing (`web-server/backend/main.go:29-31`).

## Scalability
- Horizontal: run multiple downloader instances; broker can track capacity and dispatch round-robin.
- Broker enhancements: persistence, retries, backoff, dead-letter queue.
- DB tuning: increase connections when concurrency rises (`downloader-server/src/database.go:113-115` currently conservative).
- Media storage: migrate to network shares or object storage if capacity needed.

## Known Limitations
- Frontend downloader uses a hardcoded backend host (`web-server/frontend/src/pages/Downloader.tsx:27-35`), while backend CORS allows `http://localhost:3034` (`web-server/backend/api/apis.go:201`). Align via config.
- URL support limited to YouTube.
- No authentication/authorization.
- Single active download at a time per downloader instance (queue serializes).

## Key Configuration
- `config.json` fields (`config.json:1-10`):
  - `video-storage-path`: absolute folder for videos and thumbnails.
  - `mysql-conn`, `mysql-db-name`, `mysql-table-name`: DB connectivity and schema.
  - `downloader-socket`, `message-queue-socket`, `webserver-backend-socket`: TCP endpoints.
  - `webserver-frontend-port`: frontend dev server port.

## Developer Skills Demonstrated
- Go microservices: REST API (`gorilla/mux`), TCP servers/clients, concurrency (goroutines, channels).
- Custom message queue: linked list FIFO and broker orchestration across sockets.
- Protocol design: simple text-based status and commands with lifecycle control.
- Systems integration: `yt-dlp` for acquisition; `ffmpeg` for transcoding and thumbnail extraction.
- Database engineering: connection management, table migration, conditional upsert logic, error-path cleanup.
- Configuration management: central JSON loader and shared globals across services.
- Web development: React + Vite + Tailwind, client-side routing, state management, API integration, media playback and UX polish.
- Observability and resilience: status propagation, graceful empty/missing table handling, controlled DB connection pool.
- Cross-platform considerations: MP4 transformation for Apple browser compatibility.
- LAN deployment know-how: CORS, ports, sockets tuned for local environment.

## Future Work
- Externalize frontend’s backend URL; use environment/config to avoid hardcoded IPs.
- Introduce auth and TLS for API and sockets.
- Add download concurrency with capacity-aware broker; sharded queues.
- Persist queue to disk; add retry logic and dead-letter handling.
- Implement health endpoints and readiness checks for all services.
- Extend URL support and validation beyond YouTube.

