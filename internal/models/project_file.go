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
	ProjectID uuid.UUID `gorm:"not null;type:uuid"`

	EditorID uuid.UUID `gorm:"not null;type:uuid"`
	Editor User `gorm:"foreignKey:EditorID"`

	Reviewers []User `gorm:"many2many:project_file_reviewers"`
}

func (pf *ProjectFile) BeforeCreate(tx *gorm.DB) (err error) {
	pf.ID = uuid.New()
	return err
}
