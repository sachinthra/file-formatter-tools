# Image Resizer App – Architecture & Features

## 1. Overview

A web application for resizing and compressing images/photos by:
- Dimensions (e.g., 1024x720, 500x500, etc.)
- File size (e.g., under 100KB)
- Batch processing: handle multiple images at once
- Smart cropping using image center detection (OpenCV/gocv)
- Progress feedback for large batch processing (with Redis)
- Upload images via URL
- API authentication (token-based)
- Temporary image storage using S3-compatible storage (with automatic deletion)

---

## 2. Architecture

### High-Level Diagram

```mermaid
graph TD
    A[User (Browser)] -->|Uploads, options, URLs, API Key| B[Frontend (Svelte/Vite)]
    B -->|REST API| C[Backend (Go/Gin)]
    C -->|Image Operations| D[Image Processing Libs (Go, OpenCV/gocv)]
    C -->|Progress, Jobs| E[Redis (Job/Progress Tracking)]
    C -->|Store/Serve| F[S3-compatible Storage (Temp, auto-delete)]
    C -->|Batch Output| G[Zipped Download]
    B <-->|Progress updates| C
    H[nginx] --> B
```

### Components

#### 2.1. Frontend
- **Framework:** Svelte (with Vite)
- **Served By:** nginx (Dockerized)
- **Features:**
  - Upload images (drag-and-drop, file picker, or image URLs)
  - Choose resize mode: by dimensions, max file size, or smart center crop
  - Set output format (JPG, PNG, WebP, etc.)
  - Set quality/compression level
  - Batch upload, show progress (real-time via polling or WebSockets)
  - Download all processed images (zip if batch)
  - Enter API key for authentication

#### 2.2. Backend
- **Language:** Go (Gin framework)
- **Endpoints:**
  - `/resize` (POST): Accepts images + options, returns processed images
  - `/batch` (POST): Accepts multiple images/URLs, returns zip
  - `/center-crop` (POST): Uses OpenCV/gocv for image center detection and cropping
  - `/upload-from-url` (POST): Accepts image URLs, downloads and processes them
  - `/progress/{job_id}` (GET or WS): Returns progress for ongoing batch jobs (tracked in Redis)
  - **All endpoints require API key authentication**
- **Processing:**
  - Uses disintegration/imaging, nfnt/resize, Go’s image libs, and OpenCV/gocv
  - Batch handled via goroutines
  - Job/progress info stored in **Redis**
  - Images stored in **S3** (or S3-compatible) bucket, deleted after processing/downloading
- **Authentication:**
  - Token/API-key based, checked on every request

#### 2.3. Storage
- **S3-compatible object storage** (e.g., Minio for local, AWS S3 for cloud)
- Files auto-deleted after job completion/download

#### 2.4. Progress Feedback
- **Redis** used for all job and progress tracking
- Frontend polls or subscribes for job status updates

#### 2.5. Dockerization
- **Frontend**: Svelte/Vite static build served by **nginx**
- **Backend**: Go/Gin with OpenCV, Redis client, and S3 client
- **Redis**: For job/progress tracking
- **S3-compatible storage**: For temporary image storage
- **Orchestration**: docker-compose for all components

---

## 3. Features List

### MVP Features

- [x] Upload single/multiple images (PNG, JPG, WebP, etc.)
- [x] Upload images from URL(s)
- [x] Progress feedback for large batch processing (Redis)
- [x] Resize by dimensions (width x height, maintain aspect ratio or crop)
- [x] Resize by file size (compress to fit under X KB)
- [x] Batch processing (process multiple images, return zip)
- [x] Download processed results
- [x] Smart center detection and cropping (OpenCV/gocv)
- [x] API authentication (token-based)
- [x] S3-compatible storage for images (auto-delete after job)
- [x] Simple, intuitive web UI (no login needed)
- [x] Dockerized deployment (`docker-compose up` runs everything)

### Advanced / Stretch Features

- [ ] API authentication management UI (add/remove keys)
- [ ] Support for additional formats (TIFF, BMP, etc.)
- [ ] Persistent history/logging (per session/IP)
- [ ] Error reporting/logs for failed jobs

---

## 4. Processing Logic

### Resize by Dimensions
- User selects width/height
- Option to maintain aspect ratio or crop to fit
- Resize image accordingly

### Resize by File Size
- User enters desired max file size (e.g., 100KB)
- App iteratively compresses (reducing quality/adjusting params) until under limit

### Batch Processing
- User uploads multiple images or provides multiple URLs
- Each image is processed with selected options (parallel with goroutines)
- Backend assigns unique job ID and tracks progress in **Redis**
- Return .zip file with all processed images (stored temporarily in S3, auto-deleted after download)

### Progress Feedback
- Each batch job gets a unique job ID
- Backend updates status in **Redis** (completed count, percent done, errors)
- Frontend polls `/progress/{job_id}` or uses WebSocket to update the UI

### Upload from URL
- User submits one or more image URLs
- Backend downloads images, validates format and size
- Images are processed as if they were uploaded directly

### Smart Center Detection & Cropping (OpenCV/gocv)
- Backend uses OpenCV (via gocv) to:
  1. Convert input image to grayscale
  2. Use contour or saliency detection to find the visual center (e.g. centroid of largest object)
  3. Crop the image so that the detected center is in the center of the output image (or as close as possible given aspect ratio constraints)
  4. Resize cropped result to user-specified output dimensions

---

## 5. Tech Stack Justification

- **Go (Gin):** Fast, robust, great for API server, lots of image libs
- **OpenCV/gocv:** Powerful image analysis, Go bindings
- **Svelte/Vite:** Modern, lightweight, easy to use
- **nginx:** Popular, easy, great for serving static files
- **Redis:** Reliable, fast key-value store for job/progress tracking
- **S3 (Minio or AWS):** Industry standard for object storage, scalable, easy to auto-delete
- **Docker/docker-compose:** Hassle-free local deployment, easy to extend/scale

---

## 6. Deployment

- Clone repo
- Create/configure S3 bucket (e.g., local Minio, AWS S3)
- Add API key(s) for authentication
- `docker-compose up` to build and run all services (frontend, backend, Redis, S3/minio)
- Access app at `localhost:8080` (or as mapped)

---

## 7. Example User Flows

1. **Resize by Dimension**
   - Drag-and-drop images or upload via URL, enter API key
   - Set width/height (e.g., 1024x720), maintain aspect ratio
   - Click "Process"
   - See progress bar, download processed images

2. **Resize by File Size**
   - Upload images or provide URLs, enter API key
   - Set max size = 100KB
   - Click "Process"
   - See progress, download zip

3. **Batch Processing with Progress**
   - Upload multiple images or URLs, enter API key
   - Select options (resize, file size, or smart crop)
   - Click "Process"
   - See real-time progress bar (via Redis), download zip

4. **Smart Center Crop**
   - Upload images/URLs, select “Smart Center Crop,” enter API key
   - (Optionally) set output size
   - Click “Process”
   - Download centered/cropped images

5. **Upload from URL**
   - Enter image URLs in the UI
   - Set processing options, enter API key
   - Click “Process”
   - See progress, download as with file uploads

---

## 8. Documentation Roadmap

- [ ] How to build/run (with Docker & locally)
- [ ] API docs (endpoints, options, authentication)
- [ ] Image processing concepts (resize, compression, cropping, center detection)
- [ ] Adding new features or storage/auth backends

---

## 9. License

Open source, MIT or Apache 2.0 (TBD)

---

**Confirmed decisions. Ready to proceed with initial folder structure and starter code!**