package models

import (
	"fmt"

	"github.com/google/go-github/v61/github"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProjectFile struct {
	gorm.Model

	ID uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name string `gorm:"not null"`
	Path string `gorm:"not null"`
	Version uint `gorm:"not null;default:1"`

	ProjectID uuid.UUID `gorm:"not null;type:uuid"`
	Project Project `gorm:"foreignKey:ProjectID"`

	EditorID uuid.UUID `gorm:"type:uuid"`
	Editor User `gorm:"foreignKey:EditorID"`

	CreatorID uuid.UUID `gorm:"not null;type:uuid"`
	Creator User `gorm:"foreignKey:CreatorID"`

	Reviewers []User `gorm:"many2many:project_file_reviewers"`
}

func (pf *ProjectFile) BeforeCreate(tx *gorm.DB) (err error) {
	pf.ID = uuid.New()
	return err
}

func (p Project) GetFiles(db *gorm.DB, u User) ([]ProjectFile, error) {
	var files []ProjectFile
	
	if u.Type == UserTypeWriter {
		var reviewerFiles []ProjectFile
		
		err := db.Unscoped().Table("project_file_reviewers").Joins("JOIN project_files ON project_file_reviewers.project_file_id=project_files.id").Where("project_file_reviewers.user_id = ? AND project_files.project_id = ?", u.ID, p.ID).Select("project_files.*").Preload("Creator").Preload("Editor").Preload("Reviewers").Find(&reviewerFiles).Error;

		if err != nil {
			return files, err
		}

		if err := db.Where("project_id = ? AND editor_id = ?", p.ID, u.ID).Preload("Creator").Preload("Editor").Preload("Reviewers").Find(&files).Error; err != nil {
			return files, err
		}

		return append(files, reviewerFiles...), nil
	}

	if err := db.Where("project_id = ?", p.ID).Preload("Creator").Preload("Editor").Preload("Reviewers").Find(&files).Error; err != nil {
		return files, err
	}

	return files, nil
}

func (p Project) NewFile(filename string, state *AppState) error {
	newFile := ProjectFile{
		ID: uuid.New(),
		Name: filename,
		Path: p.DirectoryPath + filename,
		Version: 1,
		ProjectID: p.ID,
		EditorID: state.User.ID,
		CreatorID: state.User.ID,
	}

	return state.DB.Transaction(func(tx *gorm.DB) error {
		gc := state.GitClient
		ctx, owner, repo := gc.Context, gc.Owner, gc.Repo

		if err := tx.Create(&newFile).Error; err != nil {
			return err
		}

		branchRef := p.GetRefName()
		commitMsg := fmt.Sprintf("%s created file %s", state.User.Username, newFile.Name)

		_, _, createErr := gc.Client.Repositories.CreateFile(ctx, owner, repo, newFile.Path, &github.RepositoryContentFileOptions{
			Message: &commitMsg,
			Branch: branchRef,
			Content: []byte{},
		})

		if createErr != nil {
			tx.Rollback()
			return createErr
		}

		return nil
	})
}

func (pf ProjectFile) GetSHA(state *AppState) (string, error) {
	gc := state.GitClient
	ctx, owner, repo := gc.Context, gc.Owner, gc.Repo

	contents, _, _, err := gc.Client.Repositories.GetContents(ctx, owner, repo, pf.Path, &github.RepositoryContentGetOptions{
		Ref: *pf.Project.GetRefName(),
	})

	if err != nil {
		return "", err
	}

	return *contents.SHA, nil
}

func (pf ProjectFile) Delete(state *AppState) error {
	gc := state.GitClient
	ctx, owner, repo := gc.Context, gc.Owner, gc.Repo
	return state.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Select(clause.Associations).Delete(&pf).Error; err != nil {
			return err
		}

		sha, err := pf.GetSHA(state)
		if err != nil {
			tx.Rollback()
			return err
		}

		commitMsg := fmt.Sprintf("%v deleted %v", state.User.Username, pf.Name)
		
		_, _, deleteErr := gc.Client.Repositories.DeleteFile(ctx, owner, repo, pf.Path, &github.RepositoryContentFileOptions{
			Message: &commitMsg,
			Branch: pf.Project.GetRefName(),
			SHA: &sha,
		})

		if deleteErr != nil {
			tx.Rollback()
			return deleteErr
		}

		return nil
	})
}
