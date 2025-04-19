package model

import (
	"time"

	"gorm.io/gorm"
)

type Folder struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	FolderName     string         `gorm:"not null;default:'/'" json:"folder_name"`
	ParentFolderID *uint          `gorm:"not null;default:0" json:"parent_folder_id"`
	UserID         uint           `gorm:"not null" json:"user_id"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	ParentFolder *Folder   `gorm:"foreignKey:ParentFolderID" json:"parent_folder"`
	SubFolders   []*Folder `gorm:"foreignKey:ParentFolderID" json:"sub_folders"`
	Files        []*File   `gorm:"foreignKey:ParentFolderID" json:"files"`
	User         *User     `gorm:"foreignKey:UserID" json:"user"`
}
