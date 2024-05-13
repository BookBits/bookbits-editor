package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccessRole string;

const (
	ReadOnlyAccessRole AccessRole = "readonly"
	ReadWriteAccessRole AccessRole = "readwrite"
)

type ProjectFileAccess struct {
	gorm.Model

	ID uuid.UUID `gorm:"primaryKey;type:uuid;"`
	ProjectFileID  uuid.UUID      `gorm:"not null"`
	ProjectFile     ProjectFile  `gorm:"foreignKey:ProjectFileID"`
	UserID        uuid.UUID      `gorm:"not null"`
	Role          AccessRole    `gorm:"not null;enum:'readonly','readwrite'"`
}

func (pfa *ProjectFileAccess) BeforeCreate(tx *gorm.DB) (err error) {
	pfa.ID = uuid.New()
	return err
}
