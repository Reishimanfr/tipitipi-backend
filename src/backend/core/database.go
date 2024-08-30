package core

import (
	"errors"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	path = "../../database.sqlite"
)

type BlogPost struct {
	ID         int `gorm:"primaryKey,autoIncrement"`
	Created_At int64
	Edited_At  int64
	Title      string `gorm:"unique"`
	Content    string
	Images     string
}

type AdminUser struct {
	ID       int `gorm:"primaryKey,autoIncrement:false"`
	Username string
	Hash     string
	Salt     string
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

	gormConfig := &gorm.Config{}

	if os.Getenv("DEV") != "true" {
		gormConfig.Logger = logger.Discard
	}

	db, err := gorm.Open(sqlite.Open(path), gormConfig)

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&BlogPost{}, &AdminUser{})
	d.DB = db

	return *d
}
