package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct{}
type Writer struct{}
type Reviewer struct{}

type UserType interface{
	UserTypeToString() string
}

func (admin Admin) UserTypeToString() string {
	return "admin"
}

func (writer Writer) UserTypeToSting() string {
	return "writer"
}

func (reviewer Reviewer) UserTypeToSting() string {
	return "reviewer"
}

type User struct {
	gorm.Model

	ID uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username string `json:"username"`
	Email string `json:"email"`
	PasswordHash string `json:"-"`
	Type UserType `json:"user_type"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return err
}
