package core

import (
	"errors"
	"os"
	"path"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Exec, _ = os.Executable()
	Path    = path.Join(Exec, "../database.sqlite")
)

type AttachmentRecord struct {
	ID         int `gorm:"primaryKey,autoIncrement"`
	Filename   string
	Path       string
	BlogPostID int `gorm:"index"`
}

type BlogPost struct {
	ID          int `gorm:"primaryKey,autoIncrement"`
	Created_At  int64
	Edited_At   int64
	Title       string `gorm:"unique"`
	Content     string
	Attachments []AttachmentRecord `gorm:"foreignKey:BlogPostID"`
}

type AdminUser struct {
	ID       int `gorm:"primaryKey,autoIncrement:false"`
	Username string
	Hash     string
	Salt     string
}

type Database struct {
	*gorm.DB
	Memory bool
}

func (d *Database) Init() Database {
	if _, err := os.Stat(Path); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(Path)

		if err != nil {
			panic(err)
		}

		defer file.Close()
	}

	gormConfig := &gorm.Config{}

	if os.Getenv("DEV") != "true" {
		gormConfig.Logger = logger.Discard
	}

	var db *gorm.DB

	if d.Memory {
		db, err = gorm.Open(sqlite.Open("sqlite::memory"), gormConfig)
	} else {
		db, err = gorm.Open(sqlite.Open(Path), gormConfig)
	}

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&BlogPost{}, &AdminUser{}, &AttachmentRecord{})
	d.DB = db

	return *d
}
