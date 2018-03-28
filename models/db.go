package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // this is needed for GORM
	"github.com/pkg/errors"
)

// DB connection
var DB *gorm.DB

// SetupDB ...
func SetupDB(connStr string) error {
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return errors.WithStack(err)
	}
	DB = db
	return nil
}

// Migrate ...
func Migrate(connStr string) error {
	if DB == nil {
		if err := SetupDB(connStr); err != nil {
			return err
		}
	}

	if err := DB.AutoMigrate(&App{}).Error; err != nil { // reason
		return err
	}
	return nil
}
