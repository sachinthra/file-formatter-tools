package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Manager struct {
	rdb *redis.Client
}

func NewManager(rdb *redis.Client) *Manager {
	return &Manager{rdb: rdb}
}

func (jm *Manager) NewJob(ctx context.Context) (string, error) {
	jobID := fmt.Sprintf("%d", time.Now().UnixNano())
	if err := jm.rdb.Set(ctx, fmt.Sprintf("job:%s:progress", jobID), 0, 30*time.Minute).Err(); err != nil {
		return "", err
	}
	return jobID, nil
}

func (jm *Manager) SetProgress(ctx context.Context, jobID string, progress int) error {
	return jm.rdb.Set(ctx, fmt.Sprintf("job:%s:progress", jobID), progress, 30*time.Minute).Err()
}

func (jm *Manager) GetProgress(jobID string) (int, error) {
	ctx := context.Background()
	val, err := jm.rdb.Get(ctx, fmt.Sprintf("job:%s:progress", jobID)).Int()
	return val, err
}

func (jm *Manager) CompleteJob(ctx context.Context, jobID string) error {
	return jm.rdb.Set(ctx, fmt.Sprintf("job:%s:progress", jobID), 100, 30*time.Minute).Err()
}
