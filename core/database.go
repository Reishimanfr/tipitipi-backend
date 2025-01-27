package core

import (
	"bash06/tipitipi-backend/flags"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PageContent struct {
	ID      int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Page    string `gorm:"unique;size:255" json:"page"`
	Content string `gorm:"type:text" json:"content"`
}

type GalleryGroup struct {
	ID          int             `gorm:"primaryKey;autoIncrement;index" json:"id"`
	Name        string          `gorm:"unique;size:255" json:"name"`
	Description string          `json:"description"`
	Images      []GalleryRecord `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"images,omitempty"`
}

type GalleryRecord struct {
	ID       int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Filename string `json:"filename"`
	GroupID  int    `gorm:"index" json:"-"`
	Mimetype string `gorm:"size:255" json:"-"`
	Size     int64  `json:"-"`
}

type File struct {
	ID         int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Filename   string `gorm:"size:255" json:"filename"`
	Size       int64  `json:"-"`
	Mimetype   string `gorm:"size:255" json:"-"`
	BlogPostID int    `gorm:"index" json:"-"`
}

type BlogPost struct {
	ID         int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Created_At int64  `gorm:"autoCreateTime" json:"created_at"`
	Edited_At  int64  `gorm:"autoUpdateTime" json:"edited_at,omitempty"`
	Title      string `gorm:"unique;size:255" json:"title"`
	Content    string `gorm:"type:text" json:"content"`
	Files      []File `gorm:"foreignKey:BlogPostID" json:"files,omitempty"`
}

type AdminUser struct {
	ID       int    `gorm:"primaryKey;autoIncrement:false"`
	Username string `gorm:"size:255"`
	Hash     string `gorm:"size:255"`
	Salt     string `gorm:"size:255"`
}

type Token struct {
	Token string `gorm:"primaryKey;size:255;index"`
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

	db.AutoMigrate(&BlogPost{}, &AdminUser{}, &File{}, &GalleryRecord{}, &GalleryGroup{}, &Token{}, PageContent{})

	return db, nil
}
