package api

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"file-formatter-tools/internal/config"
	"file-formatter-tools/internal/imgproc"
	"file-formatter-tools/internal/jobs"
	"file-formatter-tools/internal/s3"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, jobManager *jobs.Manager, s3Client *s3.Client, cfg *config.Config) {
	api := r.Group("/api")
	{
		api.GET("/progress/:jobID", ProgressHandler(jobManager))
		api.POST("/resize", ResizeHandler(s3Client, jobManager))
		api.POST("/batch", BatchHandler(s3Client, jobManager))
		api.POST("/center-crop", CenterCropHandler())
		api.POST("/upload-from-url", UploadFromURLHandler())
	}
}

// Placeholder handler: /api/resize
func ResizeHandler(s3Client *s3.Client, jobManager *jobs.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 1. Create job ID and set initial progress
		jobID, err := jobManager.NewJob(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create job"})
			return
		}
		_ = jobManager.SetProgress(ctx, jobID, 5)

		// ... (the rest of your image processing code)
		// For each major step below, update the progress:
		// e.g., jobManager.SetProgress(ctx, jobID, 20)

		// Parse form (max 32MB)
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form", "job_id": jobID})
			return
		}
		_ = jobManager.SetProgress(ctx, jobID, 10)

		// Get file
		file, header, err := c.Request.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing image file", "job_id": jobID})
			return
		}
		defer file.Close()
		_ = jobManager.SetProgress(ctx, jobID, 20)

		// Read options
		width, _ := strconv.Atoi(c.PostForm("width"))
		height, _ := strconv.Atoi(c.PostForm("height"))
		maintainAspect := c.PostForm("maintain_aspect_ratio") == "true"
		quality, _ := strconv.Atoi(c.DefaultPostForm("quality", "85"))
		maxSizeKB, _ := strconv.Atoi(c.DefaultPostForm("max_size_kb", "0")) // 0 = no limit

		// Detect extension/format
		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext == "" {
			ext = ".jpg" // default
		}

		// Read file into buffer
		imageData, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read image", "job_id": jobID})
			return
		}
		_ = jobManager.SetProgress(ctx, jobID, 30)

		// Do the resize
		result, format, err := imgproc.ResizeImage(imageData, width, height, maintainAspect, quality, maxSizeKB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Image resize failed", "details": err.Error(), "job_id": jobID})
			return
		}
		_ = jobManager.SetProgress(ctx, jobID, 60)

		// Generate object name (unique)
		uid := fmt.Sprintf("%d", time.Now().UnixNano())
		ext = strings.TrimPrefix(ext, ".")
		objectName := fmt.Sprintf("resize/%s.%s", uid, ext)
		contentType := "image/" + format

		// Upload to S3
		if err := s3Client.Upload(ctx, objectName, result, contentType); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to S3", "details": err.Error(), "job_id": jobID})
			return
		}
		_ = jobManager.SetProgress(ctx, jobID, 80)

		// Get presigned URL
		url, err := s3Client.GetPresignedURL(ctx, objectName, 10*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get download URL", "details": err.Error(), "job_id": jobID})
			return
		}
		_ = jobManager.CompleteJob(ctx, jobID)

		// Respond with download link and job ID
		c.JSON(http.StatusOK, gin.H{
			"job_id":       jobID,
			"download_url": url,
			"format":       format,
			"object_name":  objectName,
		})
	}
}

func BatchHandler(s3Client *s3.Client, jobManager *jobs.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Read options
		width, _ := strconv.Atoi(c.PostForm("width"))
		height, _ := strconv.Atoi(c.PostForm("height"))
		maintainAspect := c.PostForm("maintain_aspect_ratio") == "true"
		quality, _ := strconv.Atoi(c.DefaultPostForm("quality", "85"))
		maxSizeKB, _ := strconv.Atoi(c.DefaultPostForm("max_size_kb", "0")) // 0 = no limit

		// Parse files
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
			return
		}
		imageFiles := form.File["images"]
		if len(imageFiles) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No images provided"})
			return
		}

		// Create parent batch job
		batchJobID, err := jobManager.NewJob(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create batch job"})
			return
		}
		_ = jobManager.SetProgress(ctx, batchJobID, 0)
		numFiles := len(imageFiles)

		imageJobs := []map[string]interface{}{}
		for idx, fileHeader := range imageFiles {
			// Create sub-job for image
			jobID, _ := jobManager.NewJob(ctx)
			_ = jobManager.SetProgress(ctx, jobID, 5)

			// Open file
			file, err := fileHeader.Open()
			if err != nil {
				imageJobs = append(imageJobs, map[string]interface{}{
					"job_id":   jobID,
					"filename": fileHeader.Filename,
					"error":    "Failed to open image",
				})
				continue
			}
			imageData, err := io.ReadAll(file)
			file.Close()
			if err != nil {
				imageJobs = append(imageJobs, map[string]interface{}{
					"job_id":   jobID,
					"filename": fileHeader.Filename,
					"error":    "Failed to read image",
				})
				continue
			}
			_ = jobManager.SetProgress(ctx, jobID, 20)

			// Detect extension/format
			ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
			if ext == "" {
				ext = ".jpg" // default
			}

			// Resize
			result, format, err := imgproc.ResizeImage(imageData, width, height, maintainAspect, quality, maxSizeKB)
			if err != nil {
				imageJobs = append(imageJobs, map[string]interface{}{
					"job_id":   jobID,
					"filename": fileHeader.Filename,
					"error":    "Resize failed: " + err.Error(),
				})
				_ = jobManager.CompleteJob(ctx, jobID)
				continue
			}
			_ = jobManager.SetProgress(ctx, jobID, 60)

			// Upload to S3
			uid := fmt.Sprintf("%d", time.Now().UnixNano())
			ext = strings.TrimPrefix(ext, ".")
			objectName := fmt.Sprintf("batch/%s_%s.%s", jobID, uid, ext)
			contentType := "image/" + format
			if err := s3Client.Upload(ctx, objectName, result, contentType); err != nil {
				imageJobs = append(imageJobs, map[string]interface{}{
					"job_id":   jobID,
					"filename": fileHeader.Filename,
					"error":    "Failed to upload to S3: " + err.Error(),
				})
				_ = jobManager.CompleteJob(ctx, jobID)
				continue
			}
			_ = jobManager.SetProgress(ctx, jobID, 80)

			// Presigned URL
			url, err := s3Client.GetPresignedURL(ctx, objectName, 10*time.Minute)
			if err != nil {
				imageJobs = append(imageJobs, map[string]interface{}{
					"job_id":   jobID,
					"filename": fileHeader.Filename,
					"error":    "Failed to get download URL: " + err.Error(),
				})
				_ = jobManager.CompleteJob(ctx, jobID)
				continue
			}

			_ = jobManager.CompleteJob(ctx, jobID)

			imageJobs = append(imageJobs, map[string]interface{}{
				"job_id":       jobID,
				"filename":     fileHeader.Filename,
				"download_url": url,
				"format":       format,
				"object_name":  objectName,
			})

			// Update batch job progress
			_ = jobManager.SetProgress(ctx, batchJobID, (idx+1)*100/numFiles)
		}

		_ = jobManager.CompleteJob(ctx, batchJobID)

		c.JSON(http.StatusOK, gin.H{
			"message":      "Batch resize is experimental.",
			"batch_job_id": batchJobID,
			"image_jobs":   imageJobs,
		})
	}
}

// Placeholder handler: /api/center-crop
func CenterCropHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":  "Center-crop endpoint hit",
			"endpoint": "/api/center-crop",
		})
	}
}

// Placeholder handler: /api/upload-from-url
func UploadFromURLHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":  "Upload-from-url endpoint hit",
			"endpoint": "/api/upload-from-url",
		})
	}
}

// Handler: /api/progress/:jobID
func ProgressHandler(jobManager *jobs.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		jobID := c.Param("jobID")
		progress, err := jobManager.GetProgress(jobID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"job_id": jobID, "progress": progress})
	}
}
