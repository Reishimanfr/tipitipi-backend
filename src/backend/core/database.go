package core

import (
	"errors"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB
	err         error
	initialized = false
	path        = "./src/backend/database.sqlite"
)

type BlogPost struct {
	Id        int `gorm:"autoIncrement,primaryKey"`
	CreatedAt time.Time
	EditedAt  time.Time
	Title     string
	Content   string
	Images    string
}

func InitDatabase() error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		os.Create(path)
	}

	db, err = gorm.Open(sqlite.Open(path), &gorm.Config{})

	if err != nil {
		return err
	}

	db.AutoMigrate(&BlogPost{})

	initialized = true
	return nil
}

func GetDb() (*gorm.DB, error) {
	if !initialized {
		return nil, errors.New("database not initialized")
	}

	return db, nil
}
