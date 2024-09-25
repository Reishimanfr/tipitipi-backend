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

type GalleryRecord struct {
	ID      int    `gorm:"primaryKey;autoIncrement" json:"id"`
	AltText string `json:"alt_text,omitempty"`
	Path    string `json:"path"`
}

type AttachmentRecord struct {
	ID         int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Filename   string `json:"filename"`
	Path       string `json:"path"`
	BlogPostID int    `gorm:"index" json:"blog_post_id"`
}

type BlogPost struct {
	ID          int                `gorm:"primaryKey;autoIncrement" json:"id"`
	Created_At  int64              `json:"created_at,omitempty"`
	Edited_At   int64              `json:"edited_at,omitempty"`
	Title       string             `gorm:"unique" json:"title,omitempty"`
	Content     string             `json:"content,omitempty"`
	Attachments []AttachmentRecord `gorm:"foreignKey:BlogPostID" json:"attachments,omitempty"`
}

type AdminUser struct {
	ID       int `gorm:"primaryKey;autoIncrement:false"`
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

	// Mainly used for testing
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
