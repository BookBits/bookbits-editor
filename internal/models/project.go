package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model

	ID uuid.UUID `gorm:"primaryKey;type:uuid;"`
	Name string `gorm:"unique;not null"`
	DirectoryPath string `gorm:"unique;not null"`
	BranchName string `gorm:"unique;not null"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return err
}
