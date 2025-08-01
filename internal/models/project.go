package models

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v61/github"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Project struct {
	gorm.Model

	ID uuid.UUID `gorm:"primaryKey;type:uuid;"`
	Name string `gorm:"unique;not null;index:idx_project_name,class:FULLTEXT"`
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
		var reviewer_project_ids []uuid.UUID
		
		err := db.Table("project_file_reviewers").Joins("JOIN project_files ON project_file_reviewers.project_file_id=project_files.id").Where("project_file_reviewers.user_id = ?", u.ID).Select("project_id").Find(&reviewer_project_ids).Error;

		if err != nil {
			return projects, err
		}
		var editor_project_ids []uuid.UUID

		err = db.Preload("Editor").Model(&ProjectFile{}).Where("editor_id = ?", u.ID).Select("project_id").Find(&editor_project_ids).Error
		
		if err != nil {
			return projects, err
		}

		project_ids := append(reviewer_project_ids, editor_project_ids...)

		if err := db.Preload("Creator").Find(&projects, &project_ids).Error; err != nil {
			return projects, err
		}

		return projects, nil
	}

	if err := db.Preload("Creator").Find(&projects).Error; err != nil {
		return projects, err
	}
	return projects, nil
}

func SearchProjects(state *AppState, keyword string) ([]Project, error) {
	user := state.User
	db := state.DB

	var projects []Project

	if user.Type == UserTypeWriter {
		var reviewer_project_ids []uuid.UUID
		
		err := db.Table("project_file_reviewers").Joins("JOIN project_files ON project_file_reviewers.project_file_id=project_files.id").Where("project_file_reviewers.user_id = ?", user.ID).Select("project_id").Find(&reviewer_project_ids).Error;

		if err != nil {
			return projects, err
		}
		var editor_project_ids []uuid.UUID

		err = db.Preload("Editor").Model(&ProjectFile{}).Where("editor_id = ?", user.ID).Select("project_id").Find(&editor_project_ids).Error
		
		if err != nil {
			return projects, err
		}

		project_ids := append(reviewer_project_ids, editor_project_ids...)

		if err := db.Where("id in (?)", project_ids).Where("MATCH(name) AGAINST (? IN BOOLEAN MODE) OR name LIKE ?", keyword, keyword+"%").Find(&projects).Error; err != nil {
			return projects, err
		}

		return projects, nil
	}

	if err := db.Where("MATCH(name) AGAINST (? IN BOOLEAN MODE) OR name LIKE ?", keyword, keyword+"%").Find(&projects).Error; err != nil {
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

func (p Project) Delete(state *AppState) error {
	gc := state.GitClient
	ctx, owner, repo := gc.Context, gc.Owner, gc.Repo
	return state.DB.Transaction(func(tx *gorm.DB) error {
		var files []ProjectFile

		if err := tx.Where("project_id = ?", p.ID).Preload("Creator").Preload("Editor").Preload("Reviewers").Preload("Project").Find(&files).Error; err != nil {
			return err
		}

		for _, file := range files {
			if err := file.Delete(state); err != nil {
				return err
			}
		}

		if err := tx.Unscoped().Select(clause.Associations).Delete(&p).Error; err != nil {
			tx.Rollback()
			return err
		}

		_, delErr := gc.Client.Git.DeleteRef(ctx, owner, repo, *p.GetRefName())

		if delErr != nil {
			tx.Rollback()
			return delErr
		}
		return nil
	})
}
