# Drive Storage System - Database Documentation

## Table of Contents
1. [Overview](#overview)
2. [Database Schema](#database-schema)
3. [GORM Models](#gorm-models)
4. [Relationships](#relationships)
5. [Indexes](#indexes)
6. [Enums](#enums)
7. [Implementation Guide](#implementation-guide)

## Overview

This documentation describes the database structure for a cloud storage system with file sharing capabilities. The system is built using GORM (Go ORM) and follows a relational database design pattern.

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR,
    email VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    storage_used DOUBLE PRECISION DEFAULT 0,
    storage_limit DOUBLE PRECISION NOT NULL DEFAULT 5000000000,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Folders Table
```sql
CREATE TABLE folders (
    id SERIAL PRIMARY KEY,
    folder_name VARCHAR NOT NULL DEFAULT '/',
    parent_folder_id INTEGER REFERENCES folders(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);
```

### Files Table
```sql
CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    file_name VARCHAR NOT NULL,
    file_type VARCHAR NOT NULL,
    file_size DOUBLE PRECISION NOT NULL,
    content_type VARCHAR NOT NULL,
    folder_id INTEGER NOT NULL REFERENCES folders(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);
```

### Shares Table
```sql
CREATE TABLE shares (
    id SERIAL PRIMARY KEY,
    file_id INTEGER REFERENCES files(id),
    folder_id INTEGER REFERENCES folders(id),
    owner_id INTEGER NOT NULL REFERENCES users(id),
    shared_with_id INTEGER NOT NULL REFERENCES users(id),
    permission VARCHAR NOT NULL DEFAULT 'view',
    access_key VARCHAR NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
```

## GORM Models

### User Model
```go
type User struct {
    ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    FirstName     string    `gorm:"not null" json:"first_name"`
    LastName      string    `gorm:"default:null" json:"last_name"`
    Email         string    `gorm:"not null;unique" json:"email"`
    Password      string    `gorm:"not null" json:"-"`
    StorageUsed   float64   `gorm:"default:0" json:"storage_used"`
    StorageLimit  float64   `gorm:"not null;default:5000000000" json:"storage_limit"`
    CreatedAt     time.Time `gorm:"default:now()" json:"created_at"`
    UpdatedAt     time.Time `gorm:"default:now()" json:"updated_at"`
    Folders       []Folder  `gorm:"foreignKey:UserID" json:"folders"`
    Files         []File    `gorm:"foreignKey:UserID" json:"files"`
    SharedFiles   []Share   `gorm:"foreignKey:OwnerID" json:"shared_files"`
    ReceivedFiles []Share   `gorm:"foreignKey:SharedWithID" json:"received_files"`
}
```

### Folder Model
```go
type Folder struct {
    ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    FolderName     string    `gorm:"not null;default:/" json:"folder_name"`
    ParentFolderID *uint     `gorm:"default:null" json:"parent_folder_id"`
    UserID         uint      `gorm:"not null" json:"user_id"`
    CreatedAt      time.Time `gorm:"default:now()" json:"created_at"`
    UpdatedAt      time.Time `gorm:"default:now()" json:"updated_at"`
    DeletedAt      *time.Time `gorm:"index" json:"deleted_at"`
    
    ParentFolder   *Folder   `gorm:"foreignKey:ParentFolderID" json:"parent_folder"`
    SubFolders     []Folder  `gorm:"foreignKey:ParentFolderID" json:"sub_folders"`
    Files          []File    `gorm:"foreignKey:FolderID" json:"files"`
    User           User      `gorm:"foreignKey:UserID" json:"user"`
}
```

### File Model
```go
type FileType string

const (
    FileTypeImage    FileType = "image"
    FileTypeVideo    FileType = "video"
    FileTypePDF      FileType = "pdf"
    FileTypeDocument FileType = "document"
    FileTypeAudio    FileType = "audio"
    FileTypeOther    FileType = "other"
)

type File struct {
    ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    FileName    string    `gorm:"not null" json:"file_name"`
    FileType    FileType  `gorm:"not null" json:"file_type"`
    FileSize    float64   `gorm:"not null" json:"file_size"`
    ContentType string    `gorm:"not null" json:"content_type"`
    FolderID    uint      `gorm:"not null" json:"folder_id"`
    UserID      uint      `gorm:"not null" json:"user_id"`
    CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
    UpdatedAt   time.Time `gorm:"default:now()" json:"updated_at"`
    DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
    
    Folder      Folder    `gorm:"foreignKey:FolderID" json:"folder"`
    User        User      `gorm:"foreignKey:UserID" json:"user"`
    Shares      []Share   `gorm:"foreignKey:FileID" json:"shares"`
}
```

### Share Model
```go
type Permission string

const (
    PermissionView  Permission = "view"
    PermissionEdit  Permission = "edit"
    PermissionAdmin Permission = "admin"
)

type Share struct {
    ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
    FileID        *uint      `gorm:"default:null" json:"file_id"`
    FolderID      *uint      `gorm:"default:null" json:"folder_id"`
    OwnerID       uint       `gorm:"not null" json:"owner_id"`
    SharedWithID  uint       `gorm:"not null" json:"shared_with_id"`
    Permission    Permission `gorm:"not null;default:'view'" json:"permission"`
    AccessKey     string     `gorm:"not null" json:"access_key"`
    ExpiresAt     *time.Time `gorm:"default:null" json:"expires_at"`
    CreatedAt     time.Time  `gorm:"default:now()" json:"created_at"`
    
    File          *File      `gorm:"foreignKey:FileID" json:"file"`
    Folder        *Folder    `gorm:"foreignKey:FolderID" json:"folder"`
    Owner         User       `gorm:"foreignKey:OwnerID" json:"owner"`
    SharedWith    User       `gorm:"foreignKey:SharedWithID" json:"shared_with"`
}
```

## Relationships

### User Relationships
- One-to-Many with Folders (User has many Folders)
- One-to-Many with Files (User has many Files)
- One-to-Many with Shares (User has many Shared Files)
- One-to-Many with Received Shares (User has many Received Files)

### Folder Relationships
- Many-to-One with User (Folder belongs to User)
- Self-Referential (Folder can have Parent Folder)
- One-to-Many with Files (Folder has many Files)
- One-to-Many with SubFolders (Folder has many SubFolders)

### File Relationships
- Many-to-One with User (File belongs to User)
- Many-to-One with Folder (File belongs to Folder)
- One-to-Many with Shares (File has many Shares)

### Share Relationships
- Many-to-One with File (Share belongs to File)
- Many-to-One with Folder (Share belongs to Folder)
- Many-to-One with Owner (Share belongs to Owner User)
- Many-to-One with SharedWith (Share belongs to SharedWith User)

## Indexes

### Folder Indexes
```sql
CREATE UNIQUE INDEX idx_user_parent_name ON folders(user_id, parent_folder_id, folder_name);
CREATE INDEX idx_parent_id ON folders(parent_folder_id);
CREATE INDEX idx_user_id ON folders(user_id);
```

### File Indexes
```sql
CREATE INDEX idx_folder_id ON files(folder_id);
CREATE INDEX idx_user_id ON files(user_id);
CREATE INDEX idx_file_type ON files(file_type);
```

### Share Indexes
```sql
CREATE INDEX idx_file_id ON shares(file_id);
CREATE INDEX idx_folder_id ON shares(folder_id);
CREATE INDEX idx_owner_shared_with ON shares(owner_id, shared_with_id);
CREATE INDEX idx_access_key ON shares(access_key);
```

## Enums

### File Types
```go
type FileType string

const (
    FileTypeImage    FileType = "image"
    FileTypeVideo    FileType = "video"
    FileTypePDF      FileType = "pdf"
    FileTypeDocument FileType = "document"
    FileTypeAudio    FileType = "audio"
    FileTypeOther    FileType = "other"
)
```

### Share Permissions
```go
type Permission string

const (
    PermissionView  Permission = "view"
    PermissionEdit  Permission = "edit"
    PermissionAdmin Permission = "admin"
)
```

## Implementation Guide

### 1. Database Setup
```go
// Initialize database connection
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
    log.Fatal(err)
}

// Auto migrate tables
err = db.AutoMigrate(&User{}, &Folder{}, &File{}, &Share{})
if err != nil {
    log.Fatal(err)
}
```

### 2. Create Indexes
```go
// Create indexes after table creation
for _, model := range []interface{}{&Folder{}, &File{}, &Share{}} {
    if m, ok := model.(interface{ Indexes() []string }); ok {
        for _, idx := range m.Indexes() {
            db.Exec(fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s", idx, m.TableName()))
        }
    }
}
```

### 3. Common Operations

#### Create User
```go
user := User{
    FirstName:    "John",
    LastName:     "Doe",
    Email:        "john@example.com",
    Password:     "hashed_password",
    StorageLimit: 5000000000,
}
db.Create(&user)
```

#### Create Folder
```go
folder := Folder{
    FolderName:     "Documents",
    ParentFolderID: nil, // Root folder
    UserID:         user.ID,
}
db.Create(&folder)
```

#### Upload File
```go
file := File{
    FileName:    "document.pdf",
    FileType:    FileTypePDF,
    FileSize:    1024,
    ContentType: "application/pdf",
    FolderID:    folder.ID,
    UserID:      user.ID,
}
db.Create(&file)
```

#### Share File
```go
share := Share{
    FileID:       &file.ID,
    OwnerID:      user.ID,
    SharedWithID: recipientID,
    Permission:   PermissionView,
    AccessKey:    "unique_access_key",
}
db.Create(&share)
```

### 4. Query Examples

#### Get User's Root Folders
```go
var folders []Folder
db.Where("user_id = ? AND parent_folder_id IS NULL", userID).Find(&folders)
```

#### Get Folder Contents
```go
var files []File
db.Where("folder_id = ?", folderID).Find(&files)
```

#### Get Shared Files
```go
var shares []Share
db.Where("shared_with_id = ?", userID).Preload("File").Find(&shares)
```

### 5. Best Practices

1. **Soft Delete**: Use `DeletedAt` for soft deletes instead of hard deletes
2. **Index Usage**: Always use indexes for frequently queried fields
3. **Transaction Management**: Use transactions for multiple related operations
4. **Preloading**: Use `Preload` for related data to avoid N+1 queries
5. **Validation**: Implement input validation before database operations
6. **Error Handling**: Always check for errors after database operations

### 6. Performance Considerations

1. **Batch Operations**: Use batch operations for bulk inserts/updates
2. **Query Optimization**: Use `Select` to fetch only needed fields
3. **Pagination**: Implement pagination for large result sets
4. **Caching**: Consider caching frequently accessed data
5. **Connection Pooling**: Configure proper connection pool settings

## Conclusion

This documentation provides a comprehensive guide to the database structure and implementation of the Drive Storage System. Follow the implementation guide and best practices to ensure optimal performance and maintainability of your application. 