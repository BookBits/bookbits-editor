package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectFile struct {
	gorm.Model

	ID uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name string `gorm:"not null"`
	Path string `gorm:"not null"`
	Version uint `gorm:"not null;default:1"`
	ProjectID uuid.UUID `gorm:"not null"`
	Project Project `gorm:"foreignKey: ProjectID"`
}

func (pf *ProjectFile) BeforeCreate(tx *gorm.DB) (err error) {
	pf.ID = uuid.New()
	return err
}
