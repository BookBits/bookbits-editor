package helpers

import (
	"errors"
	"log"
	"os"

	"github.com/BookBits/bookbits-editor/internal/models"
	"github.com/joho/godotenv"
)

func SetupEnvVars() (models.EnvVars, error) {
	err := godotenv.Load()
	if err != nil {
		//log.Fatal(err)
		return models.EnvVars{}, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not provided")
		return models.EnvVars{}, errors.New("Missing Env Var")
	}

	db_port := os.Getenv("DB_PORT")
	if db_port == "" {
		log.Fatal("DB_PORT not provided")
		return models.EnvVars{}, errors.New("Missing Env Var")
	}

	db_user := os.Getenv("DB_USER")
	if db_user == "" {
		log.Fatal("DB_USER not provided")
		return models.EnvVars{}, errors.New("Missing Env Var")


	}

	db_name := os.Getenv("DB_NAME")
	if db_name == "" {
		log.Fatal("DB_NAME not provided")
		return models.EnvVars{}, errors.New("Missing Env Var")


	}
	
	db_password := os.Getenv("DB_PASSWORD")
	if db_password == "" {
		log.Fatal("DB_PASSWORD not provided")
		return models.EnvVars{}, errors.New("Missing Env Var")
	}	

	db_host := os.Getenv("DB_HOST")
	if db_host == "" {
		log.Fatal("DB_HOST not provided")
		return models.EnvVars{}, errors.New("Missing Env Var")
	}

	jwt_secret_key := os.Getenv("JWT_SECRET_KEY")
	if jwt_secret_key == "" {
		log.Fatal("JWT_SECRET_KEY not provided")
		return models.EnvVars{}, errors.New("Missing Env Var")
	}

	default_admin_user_email := os.Getenv("DEFAULT_ADMIN_USER_EMAIL")
	if default_admin_user_email == "" {
		log.Fatal("DEFAULT_ADMIN_USER_EMAIL not provided")
		return models.EnvVars{}, errors.New("Missing Env Var")
	}

	default_admin_user_password := os.Getenv("DEFAULT_ADMIN_USER_PASSWORD")
	if default_admin_user_password == "" {
		log.Fatal("DEFAULT_ADMIN_USER_PASSWORD not provided")
		return models.EnvVars{}, errors.New("Missing Env Var")
	}

	return models.EnvVars{
		Port: port,
		DbPort: db_port,
		DbHost: db_host,
		DbUser: db_user,
		DbPassword: db_password,
		DbName: db_name,
		JWTSecretKey: []byte(jwt_secret_key),
		DefaultAdminUserEmail: default_admin_user_email,
		DefaultAdminPassword: default_admin_user_password,
	}, nil
}
