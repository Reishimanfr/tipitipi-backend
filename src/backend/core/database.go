package core

import (
	"bash06/strona-fundacja/src/backend/flags"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GalleryRecord struct {
	ID       int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Filename string `json:"filename"`
	GroupID  int    `gorm:"index" json:"-"`
	Mimetype string `json:"-"`
	Size     int64  `json:"-"`
}

type GalleryGroup struct {
	ID     int             `gorm:"primaryKey;autoIncrement" json:"id"`
	Name   string          `gorm:"unique" json:"name"`
	Images []GalleryRecord `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"images,omitempty"`
}

type File struct {
	ID         int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Filename   string `json:"filename"`
	Size       int64  `json:"-"`
	Mimetype   string `json:"-"`
	BlogPostID int    `gorm:"index" json:"-"`
}

type BlogPost struct {
	ID         int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Created_At int64  `json:"created_at"`
	Edited_At  int64  `json:"edited_at,omitempty"`
	Title      string `gorm:"unique" json:"title"`
	Content    string `json:"content"`
	Files      []File `gorm:"foreignKey:BlogPostID" json:"files,omitempty"`
}

type AdminUser struct {
	ID       int `gorm:"primaryKey;autoIncrement:false"`
	Username string
	Hash     string
	Salt     string
}

func InitDb() (*gorm.DB, error) {
	gormConfig := &gorm.Config{}

	if !*flags.Dev {
		gormConfig.Logger = logger.Discard
	}

	db, err := gorm.Open(sqlite.Open("database.sqlite"), gormConfig)

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&BlogPost{}, &AdminUser{}, &File{}, &GalleryRecord{}, &GalleryGroup{})

	return db, nil
}
