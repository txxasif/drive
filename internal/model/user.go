package model

import (
	"time"

	"gorm.io/gorm"
)

// AuthProvider represents the authentication provider type
type AuthProvider string

const (
	// LocalAuth represents authentication with username/password
	LocalAuth AuthProvider = "local"
	// GoogleAuth represents authentication with Google
	GoogleAuth AuthProvider = "google"
	// FacebookAuth represents authentication with Facebook
	FacebookAuth AuthProvider = "facebook"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Email        string         `gorm:"unique;not null" json:"email"`
	Username     string         `gorm:"unique;not null" json:"username"`
	Password     string         `gorm:"not null" json:"-"`
	FirstName    string         `gorm:"not null" json:"first_name"`
	LastName     string         `json:"last_name"`
	StorageUsed  float64        `gorm:"default:0" json:"storage_used"`
	StorageLimit float64        `gorm:"default:15000"  json:"storage_limit"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	// OAuth fields
	Provider   AuthProvider `gorm:"type:varchar(20);default:'local'" json:"provider"`
	ProviderId string       `gorm:"index" json:"-"`

	Folders       []*Folder `gorm:"foreignKey:UserID" json:"folders"`
	Files         []*File   `gorm:"foreignKey:UserID" json:"files"`
	ShareFile     []*Share  `gorm:"foreignKey:OwnerId" json:"shared_file"`
	ReceivedFiles []*Share  `gorm:"foreignKey:SharedWithId" json:"received_files"`
}
