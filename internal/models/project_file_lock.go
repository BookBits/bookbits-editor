package models

import (
	"time"

	"github.com/google/uuid"
)

type ProjectFileLock struct {
	ProjectFileID uuid.UUID
	UserID uuid.UUID
	LockTime time.Time
	ExpiresAt time.Time
}
