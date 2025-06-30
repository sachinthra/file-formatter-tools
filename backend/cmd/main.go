package main

import (
	"log"
	"os"

	"file-formatter-tools/internal/api"
	"file-formatter-tools/internal/auth"
	"file-formatter-tools/internal/config"
	"file-formatter-tools/internal/jobs"
	"file-formatter-tools/internal/s3"

	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func main() {
	// Load config from env
	cfg := config.Load()

	// Set Gin mode
	if os.Getenv("GIN_MODE") != "" {
		gin.SetMode(os.Getenv("GIN_MODE"))
	}

	// Initialize Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	// Initialize S3
	s3Client := s3.NewS3Client(cfg)

	// Ensure bucket exists
	err := s3Client.CreateBucketIfNotExists(cfg.S3Bucket)
	if err != nil {
		log.Fatalf("Failed to create S3 bucket: %v", err)
	}

	// Initialize job manager
	jobManager := jobs.NewManager(rdb)

	// Gin router
	r := gin.Default()

	// Set maximum multipart memory
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	// Add CORS middleware
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:8080"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "X-API-Key"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// }))

	// Health endpoint (no auth)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API key middleware
	r.Use(auth.APIKeyAuthMiddleware(cfg.APIKeys))

	// Register routes with dependencies
	api.RegisterRoutes(r, jobManager, s3Client, cfg)

	// Debug: Print registered routes
	for _, route := range r.Routes() {
		log.Printf("Registered route: %s %s", route.Method, route.Path)
	}

	log.Printf("Backend listening on :8081 ...")
	r.Run(":8081")
}
