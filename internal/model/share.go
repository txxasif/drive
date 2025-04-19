package model

import (
	"time"

	"gorm.io/gorm"
)

type Permission string

const (
	PermissionRead  Permission = "read"
	PermissionWrite Permission = "write"
	PermissionOwner Permission = "owner"
)

type Share struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	FolderID     uint       `gorm:"not null" json:"folder_id"`
	FileID       uint       `gorm:"" json:"file_id"`
	OwnerID      uint       `gorm:"not null" json:"owner_id"`
	SharedWithID uint       `gorm:"not null" json:"shared_with_id"`
	Permission   Permission `gorm:"not null" json:"permission"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Folder     *Folder `gorm:"foreignKey:FolderID" json:"folder"`
	File       *File   `gorm:"foreignKey:FileID" json:"file"`
	Owner      *User   `gorm:"foreignKey:OwnerID" json:"owner"`
	SharedWith *User   `gorm:"foreignKey:SharedWithID" json:"shared_with"`
}
