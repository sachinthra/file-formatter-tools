package config

import (
	"os"
	"strings"
)

type Config struct {
	RedisAddr   string
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	APIKeys     []string
}

func Load() *Config {
	return &Config{
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		S3Endpoint:  getEnv("S3_ENDPOINT", "localhost:9000"),
		S3AccessKey: getEnv("S3_ACCESS_KEY", ""),
		S3SecretKey: getEnv("S3_SECRET_KEY", ""),
		S3Bucket:    getEnv("S3_BUCKET", "images"),
		APIKeys:     strings.Split(getEnv("API_KEYS", "changeme"), ","),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
