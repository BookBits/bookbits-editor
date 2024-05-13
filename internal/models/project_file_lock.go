package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectFileLock struct {
	gorm.Model

	ID uuid.UUID `gorm:"primaryKey;type:uuid"`
	ProjectFileID uuid.UUID `gorm:"not null;unique"`
	UserID uuid.UUID `gorm:"not null"`
	LockTime time.Time `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

func (lock *ProjectFileLock) BeforeCreate(tx *gorm.DB) (err error) {
	lock.ID = uuid.New()
	return err
}
