# Frontend Plan for Image Processing API (v1)

## Overview

This document describes the plan for a minimal web frontend (v1) for the image processing API.  
The frontend will enable users to:
- Upload and resize a single image.
- Upload and batch-resize multiple images.
- View job progress (for both single and batch jobs).
- Download processed results.

The frontend will be simple, user-friendly, and can be built with React (Vite), Vue, or another modern JS framework.  
(React + Vite is recommended for fast development and wide support.)

---

## Features

### 1. Single Image Resize

- **Form**: Upload a single image.
- **Options**: Width, height, maintain aspect ratio, quality, max size (KB).
- **Actions**:
  - Submit triggers `/api/resize` POST request.
  - Show job ID and link to poll progress.
  - When complete, display download link.

### 2. Batch Image Resize

- **Form**: Upload multiple images (`images[]` via multiple file select).
- **Options**: Same as above, applies to all images.
- **Actions**:
  - Submit triggers `/api/batch` POST request.
  - Show batch job ID and list of sub-job IDs.
  - Show per-image progress and download links as each result is ready.

### 3. Progress Tracking

- **Mechanism**: Poll `/api/progress/:jobID` every 1–2 seconds.
- **Display**:
  - For single jobs: Show progress bar and status.
  - For batch jobs: Show batch-level and per-image progress.

### 4. Download Results

- **UI**: Download button for each image when result is ready.
- **Batch**: Option to download all results as a .zip (future enhancement).

---

## UI Structure

- **Tabs or Sections**:
  - Single Resize
  - Batch Resize
  - (Future: Center Crop, Upload from URL, Face Crop)

- **Each Section Contains**:
  - File(s) input + options form
  - Submit button
  - Progress area
  - Download links/results

- **Notifications**: Show errors, API messages (e.g., "Batch resize is experimental.")

---

## Tech Stack

- **Frontend Framework**: React (with Vite)
- **UI Library**: (Optional) Chakra UI, Material UI, or plain CSS
- **HTTP Client**: Fetch API or Axios

---

## Example Flow

1. User selects "Single Resize" tab.
2. Uploads image, sets options, clicks "Resize."
3. UI shows job ID and progress bar.
4. When done, download link appears.

---

## Future Enhancements (v2+)

- Center crop, upload-from-URL, and face crop endpoints.
- Drag-and-drop support.
- Download all batch results as ZIP.
- Image previews.
- User authentication (if required).

---

## Directory Structure

```
frontend/
  src/
    components/
      SingleResizeForm.jsx
      BatchResizeForm.jsx
      ProgressBar.jsx
      DownloadLink.jsx
    App.jsx
    main.jsx
  public/
  index.html
  package.json
  vite.config.js
  ...
```

---

## API Reference

- **POST /api/resize**: Single image resize.
- **POST /api/batch**: Batch image resize.
- **GET /api/progress/:jobID**: Job progress.
- See backend docs for payload/response.

---

## Development & Deployment

- Use Vite for fast dev server and build.
- Optionally deploy via Netlify, Vercel, or static hosting.

---

## Example Wireframe

```
+------------------------+
| [ Tabs: Single | Batch ]|
+------------------------+
|  [Upload Form]         |
|  [Options]             |
|  [Submit Button]       |
|------------------------|
|  [Progress Bar/Status] |
|  [Download Link(s)]    |
+------------------------+
```

---

## Notes

- Keep UI simple and intuitive for v1.
- Ensure clear status/progress feedback.
- Provide error messages for failed uploads or jobs.

---