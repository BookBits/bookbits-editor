package helpers

import (
	"fmt"

	"github.com/BookBits/bookbits-editor/internal/models"
	"github.com/gofiber/fiber/v3/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDB(vars models.EnvVars) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local",
			vars.DbUser, vars.DbPassword, vars.DbHost, vars.DbPort, vars.DbName)

	log.Info(fmt.Sprintf("Connecting to DB: %v", dsn))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	migrateErr := db.AutoMigrate()
	if migrateErr != nil {
		return nil, err
	}

	return db, nil
}
