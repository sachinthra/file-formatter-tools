package s3

import (
	"bytes"
	"context"
	"log"
	"net/url"
	"time"

	"file-formatter-tools/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	Minio  *minio.Client
	Bucket string
}

func NewS3Client(cfg *config.Config) *Client {
	minioClient, err := minio.New(cfg.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		Secure: false, // Set true if using https
	})
	if err != nil {
		log.Fatalf("[FATAL] [S3] Failed to create S3 client: %v", err)
	}
	log.Printf("[INFO] [S3] Connected to S3 endpoint %s, bucket %s", cfg.S3Endpoint, cfg.S3Bucket)
	return &Client{
		Minio:  minioClient,
		Bucket: cfg.S3Bucket,
	}
}

// Uploads a file to S3, returns the object key
func (c *Client) Upload(ctx context.Context, objectName string, data []byte, contentType string) error {
	reader := bytes.NewReader(data)
	log.Printf("[INFO] [S3] Uploading object: %s (contentType=%s, size=%d)", objectName, contentType, len(data))
	_, err := c.Minio.PutObject(ctx, c.Bucket, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		log.Printf("[ERROR] [S3] Failed to upload object %s: %v", objectName, err)
	}
	return err
}

// Returns a presigned download URL, valid for "expiry" duration
func (c *Client) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)
	u, err := c.Minio.PresignedGetObject(ctx, c.Bucket, objectName, expiry, reqParams)
	if err != nil {
		log.Printf("[ERROR] [S3] Failed to get presigned URL for %s: %v", objectName, err)
		return "", err
	}
	log.Printf("[INFO] [S3] Generated presigned URL for %s", objectName)
	return u.String(), nil
}

// TODO: Implement delete methods as needed
