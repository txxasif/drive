# GORM Model Definitions

## User Model
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

## Folder Model
```go
type Folder struct {
    ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    FolderName     string    `gorm:"not null;default:/" json:"folder_name"`
    ParentFolderID *uint     `gorm:"default:null" json:"parent_folder_id"`
    UserID         uint      `gorm:"not null" json:"user_id"`
    CreatedAt      time.Time `gorm:"default:now()" json:"created_at"`
    UpdatedAt      time.Time `gorm:"default:now()" json:"updated_at"`
    DeletedAt      *time.Time `gorm:"index" json:"deleted_at"`
    
    // Relations
    ParentFolder   *Folder   `gorm:"foreignKey:ParentFolderID" json:"parent_folder"`
    SubFolders     []Folder  `gorm:"foreignKey:ParentFolderID" json:"sub_folders"`
    Files          []File    `gorm:"foreignKey:FolderID" json:"files"`
    User           User      `gorm:"foreignKey:UserID" json:"user"`
}

// Add indexes
func (f *Folder) TableName() string {
    return "folders"
}

func (f *Folder) Indexes() []string {
    return []string{
        "idx_user_parent_name", // (user_id, parent_folder_id, folder_name)
        "idx_parent_id",        // (parent_folder_id)
        "idx_user_id",          // (user_id)
    }
}
```

## File Model
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
    
    // Relations
    Folder      Folder    `gorm:"foreignKey:FolderID" json:"folder"`
    User        User      `gorm:"foreignKey:UserID" json:"user"`
    Shares      []Share   `gorm:"foreignKey:FileID" json:"shares"`
}

// Add indexes
func (f *File) TableName() string {
    return "files"
}

func (f *File) Indexes() []string {
    return []string{
        "idx_folder_id",  // (folder_id)
        "idx_user_id",    // (user_id)
        "idx_file_type",  // (file_type)
    }
}
```

## Share Model
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
    
    // Relations
    File          *File      `gorm:"foreignKey:FileID" json:"file"`
    Folder        *Folder    `gorm:"foreignKey:FolderID" json:"folder"`
    Owner         User       `gorm:"foreignKey:OwnerID" json:"owner"`
    SharedWith    User       `gorm:"foreignKey:SharedWithID" json:"shared_with"`
}

// Add indexes
func (s *Share) TableName() string {
    return "shares"
}

func (s *Share) Indexes() []string {
    return []string{
        "idx_file_id",              // (file_id)
        "idx_folder_id",            // (folder_id)
        "idx_owner_shared_with",    // (owner_id, shared_with_id)
        "idx_access_key",           // (access_key)
    }
}
```

## Usage Notes

1. Add these models to separate files in your `model` package
2. Use GORM's AutoMigrate to create tables:
```go
db.AutoMigrate(&User{}, &Folder{}, &File{}, &Share{})
```

3. For the indexes, you'll need to create them manually after table creation:
```go
for _, model := range []interface{}{&Folder{}, &File{}, &Share{}} {
    if m, ok := model.(interface{ Indexes() []string }); ok {
        for _, idx := range m.Indexes() {
            db.Exec(fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s", idx, m.TableName()))
        }
    }
}
``` 