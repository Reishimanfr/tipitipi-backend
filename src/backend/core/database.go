package core

import (
	"errors"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	err  error
	path = "../../database.sqlite"
)

type BlogPost struct {
	ID         int `gorm:"primaryKey,autoIncrement"`
	Created_At time.Time
	Edited_At  time.Time
	Title      string
	Content    string
	Images     string
}

func InitDatabase() error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(path)

		if err != nil {
			return err
		}

		defer file.Close()
	}

	db, err = gorm.Open(sqlite.Open(path), &gorm.Config{})

	if err != nil {
		return err
	}

	db.AutoMigrate(&BlogPost{})

	return nil
}

func GetDb() (*gorm.DB, error) {
	if db == nil {
		return nil, errors.New("database not initialized")
	}

	return db, nil
}
