package model

import "time"

// ExampleDTO demonstrates validation tags for a data transfer object
type ExampleDTO struct {
	// Basic validations
	Username    string `json:"username" validate:"required,min=3,max=50"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8,strong_password"`
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
	Role        string `json:"role" validate:"required,oneof=admin user guest"`

	// Number validations
	Age        int     `json:"age" validate:"required,min=18,max=120"`
	Amount     float64 `json:"amount" validate:"required,min=0.01,max=1000000"`
	Count      int     `json:"count" validate:"min=0,max=5"`
	Percentage float64 `json:"percentage" validate:"min=0,max=100"`

	// String format validations
	UUID        string `json:"uuid" validate:"uuid"`
	Website     string `json:"website" validate:"custom_url"`
	SecureURL   string `json:"secure_url" validate:"custom_url,startswith=https://"`
	DateOfBirth string `json:"date_of_birth" validate:"date"`
	AppointTime string `json:"appoint_time" validate:"time"`
	IPAddress   string `json:"ip_address" validate:"ip_address"`

	// Conditional validations
	OptionalField    string `json:"optional_field" validate:"omitempty,min=5"`
	ConditionalField string `json:"conditional_field" validate:"required_if=Role admin"`

	// Slice validations
	Tags       []string `json:"tags" validate:"required,min=1,max=5,dive,min=2,max=10"`
	Categories []string `json:"categories" validate:"omitempty,max=3,dive,min=3"`

	// Struct validations
	Address    Address    `json:"address" validate:"required"`
	CreateTime *time.Time `json:"create_time" validate:"required"`
}

// Address is a nested struct for validation example
type Address struct {
	Street  string `json:"street" validate:"required"`
	City    string `json:"city" validate:"required"`
	State   string `json:"state" validate:"required,len=2"`
	Country string `json:"country" validate:"required"`
	ZipCode string `json:"zip_code" validate:"required"`
}

/*
Validation Tags Reference:

Basic Validations:
- required: Field cannot be empty/zero
- min=n: Minimum length/value of n
- max=n: Maximum length/value of n
- len=n: Exact length/value of n
- email: Must be valid email format
- strong_password: Must have uppercase, lowercase, number, and special character

String Format Validations:
- uuid: Must be valid UUID format
- custom_url: Must be valid URL format
- date: Must be valid date format YYYY-MM-DD
- time: Must be valid time format HH:MM:SS
- ip_address: Must be valid IP address

Numeric Validations:
- min=n: Minimum value of n
- max=n: Maximum value of n
- gt=n: Greater than n
- lt=n: Less than n
- gte=n: Greater than or equal to n
- lte=n: Less than or equal to n

String Content Validations:
- startswith=x: Must start with x
- endswith=x: Must end with x
- contains=x: Must contain x
- oneof=x y z: Must be one of the provided values

Conditional Validations:
- omitempty: Only validate if field is not empty
- required_if=Field Value: Required if Field equals Value
- required_with=Field: Required if Field is not empty
- required_without=Field: Required if Field is empty

Slice/Map/Array Validations:
- dive: Validates each element with the following validations
- min=n: Minimum number of elements
- max=n: Maximum number of elements

For reference only - not intended to be used directly.
*/
