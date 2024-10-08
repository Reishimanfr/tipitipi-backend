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
	AltText string `json:"alt_text"`
	URL     string `json:"url"`
}

type AttachmentRecord struct {
	ID         int    `gorm:"primaryKey;autoIncrement" json:"id"`
	URL        string `json:"url"`
	Filename   string `json:"filename"`
	BlogPostID int    `gorm:"index" json:"-"`
}

type BlogPost struct {
	ID          int                `gorm:"primaryKey;autoIncrement" json:"id"`
	Created_At  int64              `json:"created_at"`
	Edited_At   int64              `json:"edited_at"`
	Title       string             `gorm:"unique" json:"title"`
	Content     string             `json:"content"`
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

	db, err = gorm.Open(sqlite.Open(Path), gormConfig)

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&BlogPost{}, &AdminUser{}, &AttachmentRecord{}, &GalleryRecord{})
	d.DB = db

	return *d
}
