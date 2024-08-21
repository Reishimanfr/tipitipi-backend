package core

import (
	"errors"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
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

type AdminUser struct {
	ID           int    `gorm:"primaryKey"`
	Username     string `gorm:"unique, not null"`
	PasswordHash string `gorm:"not null"`
}

type Database struct {
	*gorm.DB
}

func (d *Database) Init() Database {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(path)

		if err != nil {
			panic(err)
		}

		defer file.Close()
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&BlogPost{}, &AdminUser{})
	d.DB = db

	return *d
}
