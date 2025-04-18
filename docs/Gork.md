# GORM Model Guide

This guide provides everything you need to know about writing models using GORM in Go applications. It covers from basic to advanced concepts, with practical examples and best practices.

## Table of Contents

- [GORM Model Guide](#gorm-model-guide)
  - [Table of Contents](#table-of-contents)
  - [Getting Started](#getting-started)
    - [Installation](#installation)
    - [Basic Setup](#basic-setup)
    - [Configuration](#configuration)
  - [Basic Model Structure](#basic-model-structure)
  - [Field Types and Tags](#field-types-and-tags)
    - [Basic Field Types](#basic-field-types)
    - [Required Fields](#required-fields)
    - [Field Validation](#field-validation)
    - [Custom Field Types](#custom-field-types)
    - [Field Defaults](#field-defaults)
    - [Field Indexes](#field-indexes)
  - [Primary Keys](#primary-keys)
    - [Auto-incrementing Primary Key](#auto-incrementing-primary-key)
    - [Custom Primary Key](#custom-primary-key)
    - [Composite Primary Key](#composite-primary-key)
  - [Relationships](#relationships)
    - [One-to-One](#one-to-one)
    - [One-to-Many](#one-to-many)
    - [Many-to-Many](#many-to-many)
    - [Polymorphic](#polymorphic)
    - [Self-Referential](#self-referential)
  - [Query Building](#query-building)
    - [Basic Queries](#basic-queries)
    - [Advanced Queries](#advanced-queries)
    - [Joins](#joins)
    - [Subqueries](#subqueries)
    - [Raw SQL](#raw-sql)
  - [Transactions](#transactions)
    - [Basic Transaction](#basic-transaction)
    - [Transaction with Context](#transaction-with-context)
    - [Nested Transactions](#nested-transactions)
  - [Hooks](#hooks)
    - [Available Hooks](#available-hooks)
    - [Example Hooks](#example-hooks)
    - [Custom Hooks](#custom-hooks)
  - [Scopes](#scopes)
    - [Defining Scopes](#defining-scopes)
    - [Using Scopes](#using-scopes)
    - [Global Scopes](#global-scopes)
  - [Performance Optimization](#performance-optimization)
    - [Indexes](#indexes)
    - [Prepared Statements](#prepared-statements)
    - [Batch Operations](#batch-operations)
    - [Query Optimization](#query-optimization)
  - [Testing](#testing)
    - [Setup Test Database](#setup-test-database)
    - [Test Examples](#test-examples)
    - [Mocking](#mocking)
  - [Common Pitfalls](#common-pitfalls)
  - [Example Models](#example-models)
    - [Complete User Model](#complete-user-model)
    - [Complete Todo Model](#complete-todo-model)
    - [Complete Product Model](#complete-product-model)
  - [Request Handling and Validation](#request-handling-and-validation)
    - [Request Flow](#request-flow)
    - [Request Types](#request-types)
      - [1. Registration Request](#1-registration-request)
      - [2. Login Request](#2-login-request)
      - [3. Todo Request](#3-todo-request)
    - [Validation Process](#validation-process)
      - [1. Request Parsing](#1-request-parsing)
      - [2. Validation Rules](#2-validation-rules)
      - [3. Error Handling](#3-error-handling)
      - [4. Validation Middleware](#4-validation-middleware)
    - [Complete Request Flow Example](#complete-request-flow-example)
    - [Error Responses](#error-responses)
      - [1. Validation Errors](#1-validation-errors)
      - [2. Business Logic Errors](#2-business-logic-errors)
      - [3. Internal Server Error](#3-internal-server-error)
    - [Best Practices for Request Handling](#best-practices-for-request-handling)
    - [Relationships](#relationships-1)
    - [Query Building](#query-building-1)

## Getting Started

### Installation

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres  # or mysql, sqlite, etc.
```

### Basic Setup

```go
import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
)

func main() {
    dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
}
```

### Configuration

```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info), // Enable logging
    NamingStrategy: schema.NamingStrategy{
        TablePrefix: "t_",   // table name prefix
        SingularTable: true, // use singular table name
    },
    PrepareStmt: true, // Enable prepared statements
})
```

## Basic Model Structure

Every model should follow this basic structure:

```go
type ModelName struct {
    // Primary key
    ID        uint           `gorm:"primaryKey" json:"id"`

    // Required fields
    FieldName string         `gorm:"not null" json:"field_name"`

    // Optional fields
    OptionalField string     `json:"optional_field"`

    // Timestamps
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

## Field Types and Tags

### Basic Field Types

```go
type User struct {
    // String fields
    Name     string  `gorm:"size:100" json:"name"`
    Email    string  `gorm:"size:255;uniqueIndex" json:"email"`

    // Numeric fields
    Age      int     `gorm:"type:int" json:"age"`
    Balance  float64 `gorm:"type:decimal(10,2)" json:"balance"`

    // Boolean fields
    Active   bool    `gorm:"default:true" json:"active"`

    // Time fields
    Birthday time.Time `gorm:"type:date" json:"birthday"`

    // JSON fields
    Settings JSON     `gorm:"type:jsonb" json:"settings"`
}
```

### Required Fields

Required fields can be defined in several ways:

1. **Basic Required Field**

```go
type User struct {
    Email string `gorm:"not null" json:"email"`
}
```

2. **Required Field with Size**

```go
type User struct {
    Username string `gorm:"not null;size:50" json:"username"`
}
```

3. **Required Field with Default**

```go
type Todo struct {
    Title     string `gorm:"not null" json:"title"`
    Completed bool   `gorm:"not null;default:false" json:"completed"`
}
```

4. **Required Field with Validation**

```go
type User struct {
    Email string `gorm:"not null;uniqueIndex" json:"email" validate:"required,email"`
}
```

5. **Required Field with Custom Validation**

```go
type Todo struct {
    Title string `gorm:"not null" json:"title"`

    func (t *Todo) BeforeCreate(tx *gorm.DB) error {
        if t.Title == "" {
            return errors.New("title is required")
        }
        return nil
    }
}
```

6. **Required Field with Database Constraints**

```go
type User struct {
    Email string `gorm:"not null;uniqueIndex;size:100;check:email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$'" json:"email"`
}
```

7. **Required Field with Custom Type**

```go
type Status string

const (
    StatusActive   Status = "active"
    StatusInactive Status = "inactive"
)

type User struct {
    Status Status `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
}
```

8. **Required Field with Embedded Struct**

```go
type Address struct {
    Street  string `gorm:"not null" json:"street"`
    City    string `gorm:"not null" json:"city"`
    Country string `gorm:"not null" json:"country"`
}

type User struct {
    Address Address `gorm:"embedded" json:"address"`
}
```

9. **Required Field with Soft Delete**

```go
type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Email     string         `gorm:"not null;uniqueIndex" json:"email"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

10. **Required Field with Custom Index**

```go
type User struct {
    Email string `gorm:"not null;uniqueIndex:idx_user_email" json:"email"`
}
```

### Field Validation

Field validation in GORM can be implemented at multiple levels:

1. **Database-Level Validation**

```go
type User struct {
    // Basic not null constraint
    Email string `gorm:"not null" json:"email"`

    // Size constraint
    Username string `gorm:"size:50" json:"username"`

    // Check constraint
    Age int `gorm:"check:age >= 18" json:"age"`

    // Unique constraint
    Email string `gorm:"uniqueIndex" json:"email"`
}
```

2. **Application-Level Validation**

```go
type User struct {
    Email    string `validate:"required,email"`
    Username string `validate:"required,min=3,max=20"`
    Age      int    `validate:"required,min=18"`
}

// Custom validation function
func validatePassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    return len(password) >= 8 &&
        regexp.MustCompile(`[A-Z]`).MatchString(password) &&
        regexp.MustCompile(`[a-z]`).MatchString(password) &&
        regexp.MustCompile(`[0-9]`).MatchString(password)
}
```

3. **Hook-Based Validation**

```go
type User struct {
    Email string `gorm:"not null" json:"email"`

    func (u *User) BeforeCreate(tx *gorm.DB) error {
        if !isValidEmail(u.Email) {
            return errors.New("invalid email format")
        }
        return nil
    }
}
```

4. **Cross-Field Validation**

```go
type User struct {
    Password        string `validate:"required"`
    ConfirmPassword string `validate:"required,eqfield=Password"`

    func (u *User) BeforeCreate(tx *gorm.DB) error {
        if u.Password != u.ConfirmPassword {
            return errors.New("passwords do not match")
        }
        return nil
    }
}
```

5. **Custom Validation Types**

```go
// Custom validation type
type Email string

func (e Email) Validate() error {
    if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(string(e)) {
        return errors.New("invalid email format")
    }
    return nil
}

type User struct {
    Email Email `gorm:"not null" json:"email"`

    func (u *User) BeforeCreate(tx *gorm.DB) error {
        return u.Email.Validate()
    }
}
```

### Custom Field Types

GORM supports various custom field types:

1. **Custom String Type**

```go
type Status string

const (
    StatusActive   Status = "active"
    StatusInactive Status = "inactive"
)

type User struct {
    Status Status `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
}
```

2. **Custom JSON Type**

```go
type Address struct {
    Street  string `json:"street"`
    City    string `json:"city"`
    Country string `json:"country"`
}

type User struct {
    Address Address `gorm:"type:jsonb" json:"address"`
}
```

3. **Custom Time Type**

```go
type CustomTime time.Time

func (t *CustomTime) Scan(value interface{}) error {
    if value == nil {
        *t = CustomTime(time.Time{})
        return nil
    }
    if v, ok := value.(time.Time); ok {
        *t = CustomTime(v)
        return nil
    }
    return fmt.Errorf("cannot convert %v to CustomTime", value)
}

func (t CustomTime) Value() (driver.Value, error) {
    return time.Time(t), nil
}

type User struct {
    LastLogin CustomTime `gorm:"type:timestamp" json:"last_login"`
}
```

4. **Custom Enum Type**

```go
type UserRole string

const (
    RoleAdmin    UserRole = "admin"
    RoleUser     UserRole = "user"
    RoleManager  UserRole = "manager"
)

func (r UserRole) Valid() bool {
    switch r {
    case RoleAdmin, RoleUser, RoleManager:
        return true
    default:
        return false
    }
}

type User struct {
    Role UserRole `gorm:"type:varchar(20);not null;default:'user'" json:"role"`

    func (u *User) BeforeCreate(tx *gorm.DB) error {
        if !u.Role.Valid() {
            return errors.New("invalid role")
        }
        return nil
    }
}
```

5. **Custom Money Type**

```go
type Money struct {
    Amount   float64
    Currency string
}

func (m *Money) Scan(value interface{}) error {
    if value == nil {
        m.Amount = 0
        m.Currency = "USD"
        return nil
    }
    if v, ok := value.([]byte); ok {
        return json.Unmarshal(v, m)
    }
    return fmt.Errorf("cannot convert %v to Money", value)
}

func (m Money) Value() (driver.Value, error) {
    return json.Marshal(m)
}

type Product struct {
    Price Money `gorm:"type:jsonb" json:"price"`
}
```

6. **Custom Array Type**

```go
type Tags []string

func (t *Tags) Scan(value interface{}) error {
    if value == nil {
        *t = Tags{}
        return nil
    }
    if v, ok := value.([]byte); ok {
        return json.Unmarshal(v, t)
    }
    return fmt.Errorf("cannot convert %v to Tags", value)
}

func (t Tags) Value() (driver.Value, error) {
    return json.Marshal(t)
}

type Post struct {
    Tags Tags `gorm:"type:jsonb" json:"tags"`
}
```

7. **Custom Map Type**

```go
type Metadata map[string]interface{}

func (m *Metadata) Scan(value interface{}) error {
    if value == nil {
        *m = Metadata{}
        return nil
    }
    if v, ok := value.([]byte); ok {
        return json.Unmarshal(v, m)
    }
    return fmt.Errorf("cannot convert %v to Metadata", value)
}

func (m Metadata) Value() (driver.Value, error) {
    return json.Marshal(m)
}

type User struct {
    Metadata Metadata `gorm:"type:jsonb" json:"metadata"`
}
```

8. **Custom UUID Type**

```go
type UUID [16]byte

func (u *UUID) Scan(value interface{}) error {
    if value == nil {
        *u = UUID{}
        return nil
    }
    if v, ok := value.([]byte); ok {
        copy(u[:], v)
        return nil
    }
    return fmt.Errorf("cannot convert %v to UUID", value)
}

func (u UUID) Value() (driver.Value, error) {
    return u[:], nil
}

type User struct {
    ID UUID `gorm:"type:uuid;primaryKey" json:"id"`
}
```

9. **Custom IP Address Type**

```go
type IPAddress net.IP

func (ip *IPAddress) Scan(value interface{}) error {
    if value == nil {
        *ip = IPAddress(net.IP{})
        return nil
    }
    if v, ok := value.(string); ok {
        *ip = IPAddress(net.ParseIP(v))
        return nil
    }
    return fmt.Errorf("cannot convert %v to IPAddress", value)
}

func (ip IPAddress) Value() (driver.Value, error) {
    return ip.String(), nil
}

type User struct {
    LastIP IPAddress `gorm:"type:inet" json:"last_ip"`
}
```

10. **Custom Phone Number Type**

```go
type PhoneNumber string

func (p *PhoneNumber) Scan(value interface{}) error {
    if value == nil {
        *p = ""
        return nil
    }
    if v, ok := value.(string); ok {
        if !isValidPhoneNumber(v) {
            return errors.New("invalid phone number format")
        }
        *p = PhoneNumber(v)
        return nil
    }
    return fmt.Errorf("cannot convert %v to PhoneNumber", value)
}

func (p PhoneNumber) Value() (driver.Value, error) {
    return string(p), nil
}

type User struct {
    Phone PhoneNumber `gorm:"type:varchar(20)" json:"phone"`
}
```

### Field Defaults

1. **Static Defaults**

```go
type User struct {
    Role     string `gorm:"default:'user'" json:"role"`
    Active   bool   `gorm:"default:true" json:"active"`
    Priority int    `gorm:"default:1" json:"priority"`
}
```

2. **Dynamic Defaults**

```go
type User struct {
    CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
```

3. **Default with Function**

   ```go
   type User struct {
    func (u *User) BeforeCreate(tx *gorm.DB) error {
        if u.Role == "" {
            u.Role = "user"
        }
        return nil
    }
   }
   ```

````

### Field Indexes

1. **Basic Index**

   ```go
type User struct {
    Email string `gorm:"index" json:"email"`
}
````

2. **Unique Index**

```go
type User struct {
    Email string `gorm:"uniqueIndex" json:"email"`
}
```

3. **Composite Index**

   ```go
   type User struct {
    FirstName string `gorm:"index:idx_name" json:"first_name"`
    LastName  string `gorm:"index:idx_name" json:"last_name"`
   }
   ```

4. **Custom Index Name**

   ```go
   type User struct {
    Email string `gorm:"index:idx_user_email" json:"email"`
   }
   ```

````

5. **Full Text Index**

```go
type Post struct {
    Content string `gorm:"index:,type:gin" json:"content"`
}
````

## Primary Keys

### Auto-incrementing Primary Key

```go
type User struct {
 ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
}
```

### Custom Primary Key

```go
type User struct {
    UserID string `gorm:"primaryKey;type:uuid" json:"user_id"`
}
```

### Composite Primary Key

```go
type UserRole struct {
 UserID uint `gorm:"primaryKey" json:"user_id"`
 RoleID uint `gorm:"primaryKey" json:"role_id"`
}
```

## Relationships

### One-to-One

```go
type User struct {
    ID      uint    `gorm:"primaryKey" json:"id"`
    Profile Profile `gorm:"foreignKey:UserID" json:"profile"`
}

type Profile struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
    UserID uint   `json:"user_id"`
    User   User   `gorm:"foreignKey:UserID" json:"user"`
}
```

### One-to-Many

```go
type User struct {
    ID    uint    `gorm:"primaryKey" json:"id"`
    Posts []Post  `gorm:"foreignKey:UserID" json:"posts"`
}

type Post struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
    Title  string `json:"title"`
    UserID uint   `json:"user_id"`
    User   User   `gorm:"foreignKey:UserID" json:"user"`
}
```

### Many-to-Many

```go
type User struct {
    ID    uint   `gorm:"primaryKey" json:"id"`
    Roles []Role `gorm:"many2many:user_roles;" json:"roles"`
}

type Role struct {
    ID    uint    `gorm:"primaryKey" json:"id"`
    Users []User  `gorm:"many2many:user_roles;" json:"users"`
}
```

### Polymorphic

```go
type Comment struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Content   string    `json:"content"`
    CommentableID   uint   `json:"commentable_id"`
    CommentableType string `json:"commentable_type"`
}

type Post struct {
    ID       uint      `gorm:"primaryKey" json:"id"`
    Comments []Comment `gorm:"polymorphic:Commentable;" json:"comments"`
}

type Product struct {
    ID       uint      `gorm:"primaryKey" json:"id"`
    Comments []Comment `gorm:"polymorphic:Commentable;" json:"comments"`
}
```

### Self-Referential

```go
type Category struct {
    ID       uint       `gorm:"primaryKey" json:"id"`
    Name     string     `json:"name"`
    ParentID *uint      `json:"parent_id"`
    Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent"`
    Children []Category `gorm:"foreignKey:ParentID" json:"children"`
}
```

## Query Building

### Basic Queries

```go
// Find by primary key
db.First(&user, 1)

// Find by conditions
db.Where("name = ?", "jinzhu").First(&user)
db.Where(&User{Name: "jinzhu"}).First(&user)

// Find all records
db.Find(&users)

// Find with conditions
db.Where("name <> ?", "jinzhu").Find(&users)
db.Where("name IN ?", []string{"jinzhu", "jinzhu 2"}).Find(&users)
```

### Advanced Queries

```go
// Select specific fields
db.Select("name", "age").Find(&users)

// Order by
db.Order("age desc, name").Find(&users)

// Limit and Offset
db.Limit(10).Offset(5).Find(&users)

// Group by
db.Model(&User{}).Select("name, sum(age) as total").Group("name").Find(&results)

// Having
db.Model(&User{}).Select("name, sum(age) as total").Group("name").Having("sum(age) > ?", 100).Find(&results)
```

### Joins

```go
// Inner Join
db.Joins("Profile").Find(&users)

// Left Join
db.Joins("LEFT JOIN profiles ON profiles.user_id = users.id").Find(&users)

// Multiple Joins
db.Joins("Profile").Joins("Company").Find(&users)
```

### Subqueries

```go
// Subquery in Where
db.Where("age > ?", db.Table("users").Select("AVG(age)")).Find(&users)

// Subquery in Select
db.Select("name, (?) as total_orders",
    db.Table("orders").Select("COUNT(*)").Where("user_id = users.id"),
).Find(&users)
```

### Raw SQL

```go
// Executing Raw SQL
db.Exec("UPDATE users SET name = ? WHERE id = ?", "jinzhu", 1)

// Raw SQL with Scan
var result Result
db.Raw("SELECT name, age FROM users WHERE name = ?", "jinzhu").Scan(&result)

// Named Arguments
db.Raw("SELECT * FROM users WHERE name = @name OR age = @age",
    map[string]interface{}{"name": "jinzhu", "age": 18}).Find(&users)
```

## Transactions

### Basic Transaction

```go
tx := db.Begin()
if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}
return tx.Commit().Error
```

### Transaction with Context

```go
tx := db.WithContext(ctx).Begin()
if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}
return tx.Commit().Error
```

### Nested Transactions

```go
tx := db.Begin()
if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&profile).Error; err != nil {
        return err
    }
    return nil
}); err != nil {
    tx.Rollback()
    return err
}

return tx.Commit().Error
```

## Hooks

### Available Hooks

```go
// Creating
func (u *User) BeforeCreate(tx *gorm.DB) error
func (u *User) AfterCreate(tx *gorm.DB) error

// Updating
func (u *User) BeforeUpdate(tx *gorm.DB) error
func (u *User) AfterUpdate(tx *gorm.DB) error

// Deleting
func (u *User) BeforeDelete(tx *gorm.DB) error
func (u *User) AfterDelete(tx *gorm.DB) error

// Finding
func (u *User) AfterFind(tx *gorm.DB) error
```

### Example Hooks

```go
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // Hash password
    if u.Password != "" {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        u.Password = string(hashedPassword)
    }
    return nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
    // Send welcome email
    return nil
}
```

### Custom Hooks

```go
// Custom hook for soft delete
func (u *User) BeforeDelete(tx *gorm.DB) error {
    if u.IsAdmin {
        return errors.New("admin user cannot be deleted")
    }
    return nil
}

// Custom hook for updating timestamps
func (u *User) BeforeUpdate(tx *gorm.DB) error {
    u.UpdatedAt = time.Now()
    return nil
}
```

## Scopes

### Defining Scopes

```go
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("active = ?", true)
}

func OlderThan(age int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("age > ?", age)
    }
}
```

### Using Scopes

```go
db.Scopes(ActiveUsers, OlderThan(18)).Find(&users)
```

### Global Scopes

```go
func (u *User) ApplyGlobalScopes(db *gorm.DB) *gorm.DB {
    return db.Where("deleted_at IS NULL")
}
```

## Performance Optimization

### Indexes

```go
type User struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string `gorm:"index:idx_name"`
    Email     string `gorm:"uniqueIndex"`
    Age       int    `gorm:"index:idx_age"`
}
```

### Prepared Statements

```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    PrepareStmt: true,
})
```

### Batch Operations

```go
// Batch Insert
var users = []User{{Name: "jinzhu_1"}, {Name: "jinzhu_2"}}
db.Create(&users)

// Batch Update
db.Model(&User{}).Where("role = ?", "admin").Updates(User{Name: "jinzhu"})
```

### Query Optimization

```go
// Use Select to limit fields
db.Select("id", "name").Find(&users)

// Use Preload to avoid N+1 queries
db.Preload("Posts").Find(&users)

// Use Joins for complex queries
db.Joins("LEFT JOIN profiles ON profiles.user_id = users.id").
   Where("profiles.age > ?", 18).
   Find(&users)
```

## Testing

### Setup Test Database

```go
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatal(err)
    }

    // Migrate the schema
    db.AutoMigrate(&User{})

    return db
}
```

### Test Examples

```go
func TestUserCreate(t *testing.T) {
    db := setupTestDB(t)

    user := User{Name: "jinzhu"}
    if err := db.Create(&user).Error; err != nil {
        t.Fatal(err)
    }

    var found User
    if err := db.First(&found, user.ID).Error; err != nil {
        t.Fatal(err)
    }

    if found.Name != user.Name {
        t.Errorf("Expected name %s, got %s", user.Name, found.Name)
    }
}
```

### Mocking

```go
type MockUserRepository struct {
    CreateFunc func(ctx context.Context, user *model.User) error
    GetByIDFunc func(ctx context.Context, id uint) (*model.User, error)
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
    return m.CreateFunc(ctx, user)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
    return m.GetByIDFunc(ctx, id)
}
```

## Common Pitfalls

1. **N+1 Query Problem**

   ```go
   // Bad
   var users []User
   db.Find(&users)
   for _, user := range users {
       db.Model(&user).Association("Posts").Find(&user.Posts)
   }

   // Good
   var users []User
   db.Preload("Posts").Find(&users)
   ```

2. **Missing Error Handling**

   ```go
   // Bad
   db.Create(&user)

   // Good
   if err := db.Create(&user).Error; err != nil {
       return err
   }
   ```

3. **Incorrect Transaction Usage**

   ```go
   // Bad
   tx := db.Begin()
   tx.Create(&user)
   tx.Commit()

   // Good
   tx := db.Begin()
   if err := tx.Create(&user).Error; err != nil {
       tx.Rollback()
       return err
   }
   return tx.Commit().Error
   ```

## Example Models

### Complete User Model

```go
type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Email     string         `gorm:"uniqueIndex;not null" json:"email"`
    Username  string         `gorm:"uniqueIndex;not null" json:"username"`
    Password  string         `gorm:"not null" json:"-"`
    FirstName string         `json:"first_name"`
    LastName  string         `json:"last_name"`
    Profile   Profile        `gorm:"foreignKey:UserID" json:"profile"`
    Posts     []Post         `gorm:"foreignKey:UserID" json:"posts"`
    Roles     []Role         `gorm:"many2many:user_roles;" json:"roles"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.Password != "" {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        u.Password = string(hashedPassword)
    }
    return nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
    // Send welcome email
    return nil
}
```

### Complete Todo Model

```go
type Todo struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Title       string         `gorm:"not null;index" json:"title"`
    Description string         `json:"description"`
    Completed   bool           `gorm:"default:false" json:"completed"`
    UserID      uint           `gorm:"not null;index" json:"user_id"`
    User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Tags        []Tag          `gorm:"many2many:todo_tags;" json:"tags"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Todo) BeforeCreate(tx *gorm.DB) error {
    if t.Title == "" {
        return errors.New("title is required")
    }
    return nil
}

func (t *Todo) AfterUpdate(tx *gorm.DB) error {
    if t.Completed {
        // Send notification
    }
    return nil
}
```

### Complete Product Model

```go
type Product struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Name        string         `gorm:"not null;index" json:"name"`
    Description string         `gorm:"type:text" json:"description"`
    Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
    SKU         string         `gorm:"uniqueIndex;not null" json:"sku"`
    Stock       int           `gorm:"not null;default:0" json:"stock"`
    CategoryID  uint          `gorm:"not null;index" json:"category_id"`
    Category    Category      `gorm:"foreignKey:CategoryID" json:"category"`
    Images      []ProductImage `gorm:"foreignKey:ProductID" json:"images"`
    Reviews     []Review      `gorm:"polymorphic:Reviewable;" json:"reviews"`
    CreatedAt   time.Time     `json:"created_at"`
    UpdatedAt   time.Time     `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
    if p.SKU == "" {
        p.SKU = generateSKU()
    }
    return nil
}

func (p *Product) AfterUpdate(tx *gorm.DB) error {
    if p.Stock < 0 {
        return errors.New("stock cannot be negative")
    }
    return nil
}
```

## Request Handling and Validation

### Request Flow

1. **HTTP Request** → **Handler** → **Service** → **Repository** → **Database**
2. Each layer has specific responsibilities:
   - Handler: Request parsing, response formatting
   - Service: Business logic, validation
   - Repository: Database operations
   - Model: Data structure, validation rules

### Request Types

#### 1. Registration Request

```go
type RegisterRequest struct {
    Email     string `json:"email" validate:"required,email,max=100"`
    Username  string `json:"username" validate:"required,username"`
    Password  string `json:"password" validate:"required,password"`
    FirstName string `json:"first_name" validate:"required,name"`
    LastName  string `json:"last_name" validate:"required,name"`
}

// Validation Rules:
// - Email: Required, valid email format, max 100 chars
// - Username: Required, 3-20 chars, alphanumeric + underscore
// - Password: Required, min 8 chars, contains uppercase, lowercase, number, special char
// - FirstName/LastName: Required, 2-50 chars, letters, spaces, hyphens, apostrophes
```

#### 2. Login Request

```go
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email,max=100"`
    Password string `json:"password" validate:"required,password"`
}

// Validation Rules:
// - Email: Required, valid email format, max 100 chars
// - Password: Required, min 8 chars
```

#### 3. Todo Request

```go
type TodoCreateRequest struct {
    Title       string `json:"title" validate:"required,min=3,max=100"`
    Description string `json:"description" validate:"max=500"`
}

type TodoUpdateRequest struct {
    Title       string `json:"title" validate:"omitempty,min=3,max=100"`
    Description string `json:"description" validate:"omitempty,max=500"`
    Completed   bool   `json:"completed"`
}

// Validation Rules:
// - Title: Required for create, 3-100 chars
// - Description: Optional, max 500 chars
// - Completed: Boolean, no validation needed
```

### Validation Process

#### 1. Request Parsing

```go
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httputil.Error(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    // ... validation and processing
}
```

#### 2. Validation Rules

```go
// Custom validation functions
func validatePassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    return len(password) >= 8 &&
        regexp.MustCompile(`[A-Z]`).MatchString(password) &&
        regexp.MustCompile(`[a-z]`).MatchString(password) &&
        regexp.MustCompile(`[0-9]`).MatchString(password) &&
        regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
}

func validateUsername(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    return len(username) >= 3 && len(username) <= 20 &&
        regexp.MustCompile(`^[a-zA-Z]`).MatchString(username) &&
        regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username)
}

func validateName(fl validator.FieldLevel) bool {
    name := fl.Field().String()
    return len(name) >= 2 && len(name) <= 50 &&
        regexp.MustCompile(`^[a-zA-Z\s\-']+$`).MatchString(name) &&
        !regexp.MustCompile(`[\s\-']{2,}`).MatchString(name) &&
        !regexp.MustCompile(`^[\s\-']|[\s\-']$`).MatchString(name)
}
```

#### 3. Error Handling

```go
type ValidationError struct {
    Field      string `json:"field"`
    Message    string `json:"message"`
    StatusCode int    `json:"-"`
}

type ValidationErrors struct {
    Errors     []ValidationError `json:"errors"`
    StatusCode int              `json:"-"`
}

func (ve ValidationErrors) Error() string {
    var messages []string
    for _, err := range ve.Errors {
        messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
    }
    return strings.Join(messages, ", ")
}
```

#### 4. Validation Middleware

```go
func ValidateRequest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var req interface{}
        // Determine request type based on route
        switch r.URL.Path {
        case "/register":
            req = &RegisterRequest{}
        case "/login":
            req = &LoginRequest{}
        // ... other cases
        }

        if err := json.NewDecoder(r.Body).Decode(req); err != nil {
            httputil.Error(w, http.StatusBadRequest, "Invalid request body")
            return
        }

        if err := validation.ValidateStruct(req); err != nil {
            if ve, ok := err.(validation.ValidationErrors); ok {
                httputil.JSON(w, ve.StatusCode, ve)
                return
            }
            httputil.Error(w, http.StatusBadRequest, "Invalid request data")
            return
        }

        // Store validated request in context
        ctx := context.WithValue(r.Context(), "validated_request", req)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Complete Request Flow Example

```go
// 1. Handler
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    // Get validated request from context
    req, ok := r.Context().Value("validated_request").(*RegisterRequest)
    if !ok {
        httputil.Error(w, http.StatusBadRequest, "Invalid request")
        return
    }

    // Call service
    authResponse, err := h.authService.Register(r.Context(), req)
    if err != nil {
        switch err {
        case service.ErrUserAlreadyExists:
            httputil.Error(w, http.StatusConflict, "User already exists")
        default:
            httputil.Error(w, http.StatusInternalServerError, "Internal server error")
        }
        return
    }

    // Return response
    httputil.JSON(w, http.StatusCreated, authResponse)
}

// 2. Service
func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
    // Check if user exists
    existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err != nil {
        return nil, err
    }
    if existingUser != nil {
        return nil, ErrUserAlreadyExists
    }

    // Create user
    user := &model.User{
        Email:     req.Email,
        Username:  req.Username,
        Password:  req.Password, // Will be hashed in BeforeCreate hook
        FirstName: req.FirstName,
        LastName:  req.LastName,
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    // Generate token
    token, err := s.jwt.GenerateToken(strconv.FormatUint(uint64(user.ID), 10))
    if err != nil {
        return nil, err
    }

    return &model.AuthResponse{
        User:  user.ToResponse(),
        Token: token,
    }, nil
}

// 3. Repository
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}
```

### Error Responses

#### 1. Validation Errors

```json
{
  "errors": [
    {
      "field": "email",
      "message": "email must be a valid email address"
    },
    {
      "field": "password",
      "message": "password must be at least 8 characters long and contain uppercase, lowercase, number, and special character"
    }
  ]
}
```

#### 2. Business Logic Errors

```json
{
  "error": "User already exists"
}
```

#### 3. Internal Server Error

```json
{
  "error": "Internal server error"
}
```

### Best Practices for Request Handling

1. **Input Validation**

   - Validate at the earliest point possible
   - Use custom validation functions for complex rules
   - Provide clear error messages
   - Return appropriate HTTP status codes

2. **Error Handling**

   - Use custom error types
   - Handle errors at appropriate levels
   - Log errors with context
   - Return user-friendly error messages

3. **Security**

   - Sanitize all input
   - Use parameterized queries
   - Implement rate limiting
   - Validate content types

4. **Performance**

   - Use prepared statements
   - Implement connection pooling
   - Cache when appropriate
   - Use appropriate indexes

5. **Testing**
   - Test all validation rules
   - Test error scenarios
   - Test edge cases
   - Use mock repositories

### Relationships

GORM supports various types of relationships between models:

1. **One-to-One Relationship**

```go
type User struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string
    Profile   Profile
}

type Profile struct {
    ID        uint   `gorm:"primaryKey"`
    UserID    uint   `gorm:"uniqueIndex"`
    Bio       string
    User      User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
```

2. **One-to-Many Relationship**

```go
type User struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string
    Posts     []Post
}

type Post struct {
    ID        uint   `gorm:"primaryKey"`
    Title     string
    UserID    uint
    User      User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
```

3. **Many-to-Many Relationship**

```go
type User struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string
    Roles     []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string
    Users     []User `gorm:"many2many:user_roles;"`
}
```

4. **Polymorphic Relationship**

```go
type Comment struct {
    ID            uint   `gorm:"primaryKey"`
    Content       string
    CommentableID uint
    CommentableType string
}

type Post struct {
    ID        uint   `gorm:"primaryKey"`
    Title     string
    Comments  []Comment `gorm:"polymorphic:Commentable;"`
}

type Page struct {
    ID        uint   `gorm:"primaryKey"`
    Title     string
    Comments  []Comment `gorm:"polymorphic:Commentable;"`
}
```

5. **Self-Referential Relationship**

```go
type Category struct {
    ID            uint   `gorm:"primaryKey"`
    Name          string
    ParentID      *uint
    Parent        *Category
    Children      []Category `gorm:"foreignKey:ParentID"`
}
```

### Query Building

GORM provides a powerful query builder with various methods:

1. **Basic Queries**

```go
// Find all records
db.Find(&users)

// Find first record
db.First(&user)

// Find last record
db.Last(&user)

// Find by primary key
db.First(&user, 10)

// Find by conditions
db.Where("name = ?", "jinzhu").First(&user)
db.Where(&User{Name: "jinzhu"}).First(&user)
```

2. **Advanced Queries**

```go
// Select specific fields
db.Select("name", "age").Find(&users)

// Order by
db.Order("age desc, name").Find(&users)

// Limit and offset
db.Limit(10).Offset(5).Find(&users)

// Group by
db.Model(&User{}).Select("name, sum(age) as total").Group("name").Find(&results)

// Having
db.Model(&User{}).Select("name, sum(age) as total").Group("name").Having("sum(age) > ?", 100).Find(&results)
```

3. **Joins**

```go
// Inner join
db.Joins("Profile").Find(&users)

// Left join
db.Joins("LEFT JOIN profiles ON profiles.user_id = users.id").Find(&users)

// Multiple joins
db.Joins("Profile").Joins("Company").Find(&users)

// Join with conditions
db.Joins("Profile", db.Where(&Profile{Active: true})).Find(&users)
```

4. **Preloading**

```go
// Preload all associations
db.Preload("Profile").Find(&users)

// Preload with conditions
db.Preload("Orders", "state = ?", "paid").Find(&users)

// Nested preloading
db.Preload("Orders.OrderItems").Find(&users)

// Preload all
db.Preload(clause.Associations).Find(&users)
```

5. **Scopes**

```go
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("active = ?", true)
}

func OlderThan(age int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("age > ?", age)
    }
}

// Use scopes
db.Scopes(ActiveUsers, OlderThan(18)).Find(&users)
```

6. **Raw SQL**

```go
// Raw SQL
db.Raw("SELECT name, age FROM users WHERE name = ?", "jinzhu").Scan(&result)

// Raw SQL with variables
db.Raw("SELECT name, age FROM users WHERE name = @name", sql.Named("name", "jinzhu")).Scan(&result)

// Exec raw SQL
db.Exec("UPDATE users SET name = ? WHERE id = ?", "jinzhu", 1)
```

7. **Subqueries**

```go
// Subquery in where
db.Where("amount > (?)", db.Table("orders").Select("AVG(amount)")).Find(&orders)

// Subquery in select
db.Select("*, (?) as total", db.Table("orders").Select("SUM(amount)")).Find(&users)
```

8. **Complex Queries**

```go
// Complex query with multiple conditions
db.Where(
    db.Where("name = ?", "jinzhu").Where(db.Where("age = ?", 20).Or("age = ?", 30)),
).Or(
    db.Where("name = ?", "jinzhu2").Where("age = ?", 40),
).Find(&users)

// Query with map conditions
db.Where(map[string]interface{}{"name": "jinzhu", "age": 20}).Find(&users)

// Query with struct conditions
db.Where(&User{Name: "jinzhu", Age: 20}).Find(&users)
```

9. **Query Hooks**

```go
type User struct {
    ID   uint
    Name string
}

func (u *User) AfterFind(tx *gorm.DB) error {
    // Process after find
    return nil
}

func (u *User) BeforeQuery(tx *gorm.DB) error {
    // Process before query
    return nil
}
```

10. **Query Optimization**

```go
// Use index hints
db.Clauses(hints.UseIndex("idx_user_name")).Find(&users)

// Use force index
db.Clauses(hints.ForceIndex("idx_user_name")).Find(&users)

// Use ignore index
db.Clauses(hints.IgnoreIndex("idx_user_name")).Find(&users)

// Use straight join
db.Clauses(hints.StraightJoin()).Find(&users)
```
