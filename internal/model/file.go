package model

import (
	"time"

	"gorm.io/gorm"
)

type FileType string

const (
	FileTypeImage FileType = "image"
	FileTypeVideo FileType = "video"
	FileTypeAudio FileType = "audio"
	FileTypePDF   FileType = "document"
	FileTypeOther FileType = "other"
)

type File struct {
	ID       uint     `gorm:"primaryKey" json:"id"`
	FileName string   `gorm:"not null" json:"file_name"`
	FileType FileType `gorm:"not null" json:"file_type"`
	FileSize int64    `gorm:"not null" json:"file_size"`
	FileURL  string   `gorm:"not null" json:"file_url"`
	FolderID uint     `gorm:"not null" json:"folder_id"`
	UserID   uint     `gorm:"not null" json:"user_id"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Folder *Folder `gorm:"foreignKey:FolderID" json:"folder"`
	User   *User   `gorm:"foreignKey:UserID" json:"user"`
	Shares []Share `gorm:"many2many:file_shares;" json:"shares"`
}
