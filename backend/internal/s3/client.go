package s3

import (
	"bytes"
	"context"
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
		panic(err)
	}
	return &Client{
		Minio:  minioClient,
		Bucket: cfg.S3Bucket,
	}
}

// Uploads a file to S3, returns the object key
func (c *Client) Upload(ctx context.Context, objectName string, data []byte, contentType string) error {
	reader := bytes.NewReader(data)
	_, err := c.Minio.PutObject(ctx, c.Bucket, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// Returns a presigned download URL, valid for "expiry" duration
func (c *Client) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)
	u, err := c.Minio.PresignedGetObject(ctx, c.Bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// TODO: Implement delete methods as needed
