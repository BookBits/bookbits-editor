package models

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (writer Writer) UserTypeToString() string {
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
	PasswordHash []byte `json:"-"`
	Type UserType `json:"user_type"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return err
}

func GetUserByID(id uuid.UUID, db *gorm.DB) (User, error) {
	var user User;
	err := db.First(&user, id).Error
	return user, err
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
	  log.Fatal(err)
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
