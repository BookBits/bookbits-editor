package models

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserType string

const (
  UserTypeAdmin   UserType = "admin"
  UserTypeReviewer UserType = "reviewer"
  UserTypeWriter  UserType = "writer"
)

func (ut UserType) ToString() string {
	switch ut {
	case UserTypeAdmin:
		return "Admin"
	case UserTypeReviewer:
		return "Reviewer"
	case UserTypeWriter:
		return "Writer"
	}
	return ""
}

type User struct {
	gorm.Model

	ID uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;"`
	Username string `json:"username"`
	Email string `json:"email" gorm:"unique"`
	PasswordHash []byte `json:"-"`
	PasswordSalt []byte `json:"-"`
	Type UserType `json:"user_type"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return err
}

func (u User) UpdatePassword(current string, newPass string, db *gorm.DB) error {
	validateErr := u.Validate(current)
	if validateErr != nil {
		return fmt.Errorf("Current Password provided is incorrect")
	}

	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}
	hasedPassword, err := bcrypt.GenerateFromPassword(append(salt, []byte(newPass)...), 12)

	if err != nil {
		return fmt.Errorf("Couldn't complete request due to server error")
	}

	saveErr := db.Model(&u).Updates(&User{PasswordSalt: salt, PasswordHash: hasedPassword}).Error
	
	if saveErr != nil {
		return fmt.Errorf("Couldn't complete request due to server error")
	}
	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generatePassword(length int) (string, error) {
    buffer := make([]byte, length)
    _, err := rand.Read(buffer)
    if err != nil {
        return "", fmt.Errorf("error generating random bytes: %w", err)
    }
    for i := range buffer {
        buffer[i] = letterBytes[int(buffer[i])%len(letterBytes)]
    }
    return string(buffer), nil
}

func (u User) UpdatePasswordRandom(db *gorm.DB) (string, error) {
	newPass, err := generatePassword(12)
	if err != nil {
		return "", err
	}
	salt := make([]byte, 16)
	_, randErr := rand.Read(salt)
	if randErr != nil {
		return "", err
	}
	hasedPassword, err := bcrypt.GenerateFromPassword(append(salt, []byte(newPass)...), 12)

	if err != nil {
		return "", fmt.Errorf("Couldn't complete request due to server error")
	}

	saveErr := db.Model(&u).Updates(&User{PasswordSalt: salt, PasswordHash: hasedPassword}).Error
	
	if saveErr != nil {
		return "", fmt.Errorf("Couldn't complete request due to server error")
	}
	return newPass, nil
}

func GetUserByID(id uuid.UUID, db *gorm.DB) (User, error) {
	var user User;
	err := db.First(&user, id).Error
	return user, err
}

func GetUserByEmail(email string, db *gorm.DB) (User, error) {
	var user User
	err := db.Where("email = ?", email).Find(&user).Error
	return user, err
}

func GetUsers(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Find(&users).Error
	return users, err
}

func CreateUserWithPassword(username string, email string, password string, user_type UserType, db *gorm.DB) error {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}
	hasedPassword, err := bcrypt.GenerateFromPassword(append(salt, []byte(password)...), 12)
	if err != nil {
		return err
	}

	user := User{
		ID: uuid.New(),
		Username: username,
		Email: email,
		PasswordHash: hasedPassword,
		PasswordSalt: salt,
		Type: user_type,
	}

	createErr := db.Create(&user).Error
	return createErr
}

func DeleteUserByID(id uuid.UUID, db *gorm.DB) error {
	err := db.Unscoped().Delete(&User{}, id).Error
	return err
}

func (user User) Validate(password string) error {
	return bcrypt.CompareHashAndPassword(user.PasswordHash, append(user.PasswordSalt, []byte(password)...))
}

type UserClaims struct {
	UserID uuid.UUID
	Username string
	UserEmail string

	jwt.RegisteredClaims
}

func (user User) GenerateTokens(vars EnvVars) (string, string, error) {
	claims := UserClaims{
		UserID: user.ID,
		Username: user.Username,
		UserEmail: user.Email,
	}

	key := vars.JWTSecretKey

	accessToken, err := generateAccessToken(claims, key)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateRefreshToken(claims, key)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ValidateToken(tokenString string, key []byte) (*UserClaims, error) {
  token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {return key, nil})
  if err != nil {
    if errors.Is(err, jwt.ErrSignatureInvalid) {
      return nil, errors.New("invalid token signature")
    } else {
      return nil, err
    }
  }

  if !token.Valid {
    return nil, errors.New("invalid token")
  }

  claims, ok := token.Claims.(*UserClaims)
  if !ok {
    return nil, errors.New("invalid token claims")
  }

  return claims, nil
}

func generateJWT(claims UserClaims, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func generateAccessToken(claims UserClaims, key []byte) (string, error) {
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 1)) // Set expiry for 1 hour
	return generateJWT(claims, key)
}

func generateRefreshToken(claims UserClaims, key []byte) (string, error) {
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)) // Set expiry for 1 week
	return generateJWT(claims, key)
}
