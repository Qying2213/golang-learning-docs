package database

import (
	"log"

	"github.com/qinyang/taskmanager/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(dsn string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	log.Println("Database connected successfully")
	return AutoMigrate()
}

func AutoMigrate() error {
	log.Println("Running database migrations...")
	err := DB.AutoMigrate(
		&models.User{},
		&models.Task{},
	)
	if err != nil {
		return err
	}
	log.Println("Database migrations completed")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
