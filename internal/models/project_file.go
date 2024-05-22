package models

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/gofiber/fiber/v3/log"
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

func (file ProjectFile) GetContentFromGit(state *AppState) (string, error) {
	gc := state.GitClient
	ctx, owner, repo := gc.Context, gc.Owner, gc.Repo

	brachRef := file.Project.GetRefName()

	contents, _, _, err := gc.Client.Repositories.GetContents(ctx, owner, repo, file.Path, &github.RepositoryContentGetOptions{
		Ref: *brachRef,
	})

	if err != nil {
		return "", err
	}

	decodedContent, err := base64.StdEncoding.DecodeString(*contents.Content)

	return string(decodedContent), err
}

func (file ProjectFile) GetContents(state *AppState) (string, error) {
	appCache := state.Cache
	ctx := context.TODO()
	cacheID := fmt.Sprintf("file-contents-%v-%v", file.ID, file.Version)

	var cachedContents string
	err := appCache.Get(ctx, cacheID, cachedContents)

	if err != nil {
		contents, err := file.GetContentFromGit(state)
		if err != nil {
			return "", err
		}

		appCache.Set(&cache.Item{
			Key: cacheID,
			Value: contents,
			TTL: time.Hour * 24 * 7,
		})

		return contents, nil
	}

	return cachedContents, nil
}

func (file ProjectFile) Save(state *AppState, contents []byte) (uint, error) {
	db := state.DB
	appCache := state.Cache
	cacheCtx := context.TODO()
	cacheID := fmt.Sprintf("file-contents-%v-%v", file.ID, file.Version)
	gc := state.GitClient
	ctx, owner, repo := gc.Context, gc.Owner, gc.Repo

	delErr := appCache.Delete(cacheCtx, cacheID)
	if delErr != nil {
		return file.Version, delErr
	}

	newVersion := file.Version + 1

	if err := db.Model(&file).Update("version", newVersion).Error; err != nil {
		return file.Version, err
	}
	
	cacheID = fmt.Sprintf("file-contents-%v-%v", file.ID, newVersion)
	cacheSaveErr := appCache.Set(&cache.Item{
		Key: cacheID,
		Value: string(contents),
		TTL: time.Hour * 24 * 7,
	})

	if cacheSaveErr != nil {
		return file.Version, cacheSaveErr
	}

	go func() {
		for try := 1; try <= 3; try++ {
			commitMsg := fmt.Sprintf("user %v saved changes on %v", state.User.Username, time.Now())
			sha, shaErr := file.GetSHA(state)

			if shaErr != nil {
				continue
			}

			_, _, gitSaveErr := gc.Client.Repositories.UpdateFile(ctx, owner, repo, file.Path, &github.RepositoryContentFileOptions{
				Message: &commitMsg,
				Branch: &file.Project.BranchName,
				SHA: &sha,
				Content: contents,
			})

			if gitSaveErr == nil {
				log.Infof("file %v(%v) saved successfully", file.ID, file.Name)
				break
			}

			//notify admins of save sync with git failed
		}
	}()

	return newVersion, nil
}
