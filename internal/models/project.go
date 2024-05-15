package models

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v61/github"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model

	ID uuid.UUID `gorm:"primaryKey;type:uuid;"`
	Name string `gorm:"unique;not null"`
	DirectoryPath string `gorm:"unique;not null"`
	BranchName string `gorm:"unique;not null"`

	CreatorID uuid.UUID `gorm:"not null;type:uuid;"`
	Creator User `gorm:"foreignKey:CreatorID"`
}

func (p Project) GetRefName() *string {
	ref := fmt.Sprintf("refs/heads/%s", p.BranchName)
	return &ref
}

func (p *Project) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return err
}

func (u User) GetProjects(db *gorm.DB) ([]Project, error) {
	var projects []Project

	if u.Type == UserTypeWriter {
		return projects, nil
	}

	if err := db.Preload("Creator").Find(&projects).Error; err != nil {
		return projects, err
	}
	return projects, nil
}

func NewProject(name string, db *gorm.DB, gc GitClient, creatorID uuid.UUID) error {
	proj := Project{
		ID: uuid.New(),
		Name: name,
		DirectoryPath: fmt.Sprintf("project-%s/", strings.ReplaceAll(strings.ToLower(name), " ", "-")),
		BranchName:  fmt.Sprintf("project-%s", strings.ReplaceAll(strings.ToLower(name), " ", "-")),
		CreatorID: creatorID,
	}

	return db.Transaction(func(tx *gorm.DB) error {
		ctx, owner, repo := gc.Context, gc.Owner, gc.Repo

		if err := tx.Create(&proj).Error; err != nil {
			return err
		}

		mainBranchRef, _, err := gc.Client.Git.GetRef(ctx, owner, repo, "heads/main")
		if err != nil {
			tx.Rollback()
			return err
		}

		newBranchRef := &github.Reference{
			Ref: proj.GetRefName(),
			Object: mainBranchRef.Object,
		}

		ref, _, err := gc.Client.Git.CreateRef(ctx, owner, repo, newBranchRef)

		if err != nil {
			tx.Rollback()
			return err
		}

		readmePath := proj.DirectoryPath + "README.md"
		readmeContents := []byte("Project: " + proj.Name)

		commitMsg := fmt.Sprintf("Initial Commit for Project %s", proj.Name)
		_, _, commitErr := gc.Client.Repositories.CreateFile(ctx, owner, repo, readmePath, &github.RepositoryContentFileOptions{
			Message: &commitMsg,
			Content: readmeContents,
			Branch: ref.Ref,
		})

		if commitErr != nil {
			tx.Rollback()
			return commitErr
		}

		return nil
	})
}
