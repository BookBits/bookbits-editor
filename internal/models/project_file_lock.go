package models

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/google/uuid"
)

type ProjectFileLock struct {
	ProjectFileID uuid.UUID
	UserID uuid.UUID
	LockTime time.Time
	ExpiresAt time.Time
}

func (file ProjectFile) LockFile(state *AppState) (ProjectFileLock, error) {
	appCache := state.Cache
	user := state.User
	lockID := fmt.Sprintf("file-lock-%v", file.ID)
	lockTime := time.Now()
	expireDuration := time.Minute * 30
	expiresAt := lockTime.Add(expireDuration)

	lock := ProjectFileLock{
		ProjectFileID: file.ID,
		UserID: user.ID,
		LockTime: lockTime,
		ExpiresAt: expiresAt,
	}

	err := appCache.Set(&cache.Item{
		Key: lockID,
		Value: lock,
		TTL: expireDuration,
	})

	return lock, err
} 

func (file ProjectFile) IsLocked(state *AppState) (uuid.UUID, error) {
	appCache := state.Cache
	ctx := context.TODO()
	lockID := fmt.Sprintf("file-lock-%v", file.ID)

	exists := appCache.Exists(ctx, lockID)

	if !exists {
		return uuid.Nil, nil
	}

	var lock ProjectFileLock
	err := appCache.Get(ctx, lockID, &lock)

	return lock.UserID, err
}
