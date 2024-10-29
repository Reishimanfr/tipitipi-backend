package core

import (
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
	Key     string `json:"key"` // AWS bucket key
	GroupID int    `gorm:"index" json:"group_id"`
}

type GalleryGroup struct {
	ID     int             `gorm:"primaryKey;autoIncrement" json:"id"`
	Name   string          `gorm:"unique" json:"name"`
	Images []GalleryRecord `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"images"`
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
	Edited_At   int64              `json:"edited_at,omitempty"`
	Title       string             `gorm:"unique" json:"title"`
	Content     string             `json:"content,omitempty"`
	Attachments []AttachmentRecord `gorm:"foreignKey:BlogPostID" json:"attachments,omitempty"`
}

type AdminUser struct {
	ID       int `gorm:"primaryKey;autoIncrement:false"`
	Username string
	Hash     string
	Salt     string
}

func InitDb(testing bool) (*gorm.DB, error) {
	gormConfig := &gorm.Config{}

	if os.Getenv("DEV") != "true" {
		gormConfig.Logger = logger.Discard
	}

	var db *gorm.DB

	if testing {
		db, err = gorm.Open(sqlite.Open("file::memory:"))
	} else {
		db, err = gorm.Open(sqlite.Open(Path), gormConfig)
	}

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&BlogPost{}, &AdminUser{}, &AttachmentRecord{}, &GalleryRecord{}, &GalleryGroup{})

	return db, nil
}
