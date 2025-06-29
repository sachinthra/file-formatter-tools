package main

import (
	"log"
	"os"

	"file-formatter-tools/internal/api"
	"file-formatter-tools/internal/auth"
	"file-formatter-tools/internal/config"
	"file-formatter-tools/internal/jobs"
	"file-formatter-tools/internal/s3"

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

	// Initialize job manager
	jobManager := jobs.NewManager(rdb)

	// Gin router
	r := gin.Default()

	// API key middleware
	r.Use(auth.APIKeyAuthMiddleware(cfg.APIKeys))

	// Register routes with dependencies
	api.RegisterRoutes(r, jobManager, s3Client, cfg)

	// Health endpoint (no auth)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Printf("Backend listening on :8081 ...")
	r.Run(":8081")
}
