package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Manager struct {
	rdb *redis.Client
}

func NewManager(rdb *redis.Client) *Manager {
	log.Printf("[INFO] [Jobs] Redis job manager initialized")
	return &Manager{rdb: rdb}
}

func (jm *Manager) NewJob(ctx context.Context) (string, error) {
	jobID := fmt.Sprintf("%d", time.Now().UnixNano())
	if err := jm.rdb.Set(ctx, fmt.Sprintf("job:%s:progress", jobID), 0, 30*time.Minute).Err(); err != nil {
		log.Printf("[ERROR] [Jobs] Failed to create new job: %v", err)
		return "", err
	}
	log.Printf("[INFO] [Jobs] Created new job: %s", jobID)
	return jobID, nil
}

func (jm *Manager) SetProgress(ctx context.Context, jobID string, progress int) error {
	err := jm.rdb.Set(ctx, fmt.Sprintf("job:%s:progress", jobID), progress, 30*time.Minute).Err()
	if err != nil {
		log.Printf("[ERROR] [Jobs] Failed to set progress for job %s: %v", jobID, err)
	}
	return err
}

func (jm *Manager) GetProgress(jobID string) (int, error) {
	ctx := context.Background()
	val, err := jm.rdb.Get(ctx, fmt.Sprintf("job:%s:progress", jobID)).Int()
	if err != nil {
		log.Printf("[ERROR] [Jobs] Failed to get progress for job %s: %v", jobID, err)
	}
	return val, err
}

func (jm *Manager) CompleteJob(ctx context.Context, jobID string) error {
	err := jm.rdb.Set(ctx, fmt.Sprintf("job:%s:progress", jobID), 100, 30*time.Minute).Err()
	if err != nil {
		log.Printf("[ERROR] [Jobs] Failed to complete job %s: %v", jobID, err)
	}
	return err
}
