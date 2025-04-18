# Complete Guide to the Simplified Validation System

This document provides an in-depth explanation of how the simplified validation system works in our application, from request validation to error responses.

## Table of Contents

1. [Overview](#overview)
2. [Key Components](#key-components)
3. [How Validation Works](#how-validation-works)
4. [Validation Tags](#validation-tags)
5. [Custom Validators](#custom-validators)
6. [Handler Integration](#handler-integration)
7. [Error Responses](#error-responses)
8. [Best Practices](#best-practices)
9. [Examples](#examples)
10. [Technical Details](#technical-details)

## Overview

Our validation system is built on top of the popular [go-playground/validator](https://github.com/go-playground/validator) package, with additional abstractions to make it more developer-friendly and integrate smoothly with our API response system.

The primary goals of this system are:

- **Declarative Validation**: Define validation rules directly in your structs using tags
- **Field-Based Errors**: Generate clear error messages specific to each field
- **Simple Integration**: Easy to use in HTTP handlers with minimal code
- **Consistent Responses**: Format errors in a standardized way for frontend consumption
- **Performance Optimization**: Pre-compiled regular expressions for efficient validation

## Key Components

### 1. The Validator Instance

```go
var (
    validate *validator.Validate
    // Pre-compile all regular expressions
    urlRegex  = regexp.MustCompile(`^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(?:/[a-zA-Z0-9\-\._~:/?#[\]@!$&'()*+,;=]*)?$`)
    dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
    timeRegex = regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d):([0-5]\d)$`)
    ipv4Regex = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
    ipv6Regex = regexp.MustCompile(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`)
)

func init() {
    validate = validator.New()

    // Configure to use JSON field names
    validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
        name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
        if name == "-" {
            return fld.Name
        }
        return name
    })

    // Register custom validators
    validate.RegisterValidation("strong_password", strongPassword)
    validate.RegisterValidation("custom_url", isURL)
    validate.RegisterValidation("date", isDate)
    validate.RegisterValidation("time", isTime)
    validate.RegisterValidation("ip_address", isIPAddress)
}
```

### 2. Validation Functions

The following key functions handle validation:

- `ValidateStructWithFields(i interface{}) map[string]string` - Validates a struct and returns field-mapped errors
- `ValidateRequestWithFields(r *http.Request, dst interface{}) map[string]string` - Decodes and validates with field-mapped errors

### 3. Error Formatting

Error messages are formatted in a user-friendly way:

```go
// formatValidationErrorMessage returns just the message part of a validation error
func formatValidationErrorMessage(err validator.FieldError) string {
    switch err.Tag() {
    case "required":
        return "is required"
    case "email":
        return "must be a valid email address"
    case "min":
        return "must be at least " + err.Param() + " characters long"
    case "max":
        return "must be at most " + err.Param() + " characters long"
    case "strong_password":
        return "must contain at least 1 uppercase, 1 lowercase, 1 number, and 1 special character"
    case "custom_url":
        return "must be a valid URL"
    case "date":
        return "must be a valid date in format YYYY-MM-DD"
    case "time":
        return "must be a valid time in format HH:MM:SS"
    case "ip_address":
        return "must be a valid IP address"
    default:
        return "failed validation: " + err.Tag()
    }
}
```

### 4. Response Utilities

The response package provides utilities to send validation errors:

- `ValidationError(w http.ResponseWriter, details ...string)` - Sends array-based validation errors
- `ValidationErrorWithFields(w http.ResponseWriter, fields map[string]string)` - Sends field-mapped validation errors

## How Validation Works

The validation process follows these steps:

1. **Request Decoding**: The JSON request body is decoded into your struct
2. **Struct Validation**: Each field is checked against its validation tags
3. **Error Collection**: Any validation errors are collected and formatted
4. **Response Generation**: Errors are returned in a standardized JSON structure

### Behind the Scenes

When you call `ValidateRequestWithFields()`:

1. The function decodes the request body using `json.NewDecoder(r.Body).Decode(dst)`
2. It then calls `validate.Struct(dst)` to validate against the struct's tags
3. If there are errors, it builds a map of field name -> error message
4. Your handler sends this map using `response.ValidationErrorWithFields()`

This is the preferred method for validating requests in handlers, as it provides field-specific error information which creates a better user experience.

## Validation Tags

### Basic Validations

```go
type User struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"min=18,max=120"`
}
```

- `required`: Field cannot be empty/zero
- `min=n`: Minimum length (for strings) or value (for numbers)
- `max=n`: Maximum length (for strings) or value (for numbers)
- `len=n`: Exact length/value
- `email`: Must be valid email format
- `strong_password`: Must have uppercase, lowercase, number, and special characters

### Format Validations

```go
type Profile struct {
    PersonalID string `json:"personal_id" validate:"uuid"`
    Website    string `json:"website" validate:"custom_url"`
    BirthDate  string `json:"birth_date" validate:"date"`
    MeetingTime string `json:"meeting_time" validate:"time"`
    ServerIP   string `json:"server_ip" validate:"ip_address"`
}
```

- `uuid`: Must be valid UUID format
- `custom_url`: Must be valid URL format (uses pre-compiled regex pattern)
- `date`: Must be valid date in YYYY-MM-DD format
- `time`: Must be valid time in HH:MM:SS format
- `ip_address`: Must be valid IP address (supports both IPv4 and IPv6)

### String Content Validations

```go
type Configuration struct {
    Protocol  string `json:"protocol" validate:"oneof=http https"`
    Prefix    string `json:"prefix" validate:"startswith=api-"`
    Suffix    string `json:"suffix" validate:"endswith=.json"`
    Contains  string `json:"contains" validate:"contains=configuration"`
}
```

- `startswith=x`: Must start with x
- `endswith=x`: Must end with x
- `contains=x`: Must contain x
- `oneof=x y z`: Must be one of the provided values

### Conditional Validations

```go
type Payment struct {
    Method      string  `json:"method" validate:"required,oneof=credit paypal bank"`
    CardNumber  string  `json:"card_number" validate:"omitempty,min=16,max=16"`
}
```

- `omitempty`: Only validate if field is not empty

## Custom Validators

### How Custom Validators Work

Custom validators are functions that implement the `validator.Func` type:

```go
type Func func(fl FieldLevel) bool
```

For example, here's how the `strong_password` validator is implemented:

```go
// strongPassword validates if a password is strong
func strongPassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()

    var hasUpper, hasLower, hasNumber, hasSpecial bool

    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsNumber(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }

        // Early return if all criteria are met
        if hasUpper && hasLower && hasNumber && hasSpecial {
            return true
        }
    }

    return hasUpper && hasLower && hasNumber && hasSpecial
}
```

### Other Custom Validators

```go
// isURL checks if a string is a valid URL
func isURL(fl validator.FieldLevel) bool {
    return urlRegex.MatchString(fl.Field().String())
}

// isDate checks if a string is a valid date in YYYY-MM-DD format
func isDate(fl validator.FieldLevel) bool {
    return dateRegex.MatchString(fl.Field().String())
}

// isTime checks if a string is a valid time in HH:MM:SS format
func isTime(fl validator.FieldLevel) bool {
    return timeRegex.MatchString(fl.Field().String())
}

// isIPAddress checks if a string is a valid IP address (IPv4 or IPv6)
func isIPAddress(fl validator.FieldLevel) bool {
    ip := fl.Field().String()
    return ipv4Regex.MatchString(ip) || ipv6Regex.MatchString(ip)
}
```

### Registering Custom Validators

Custom validators are registered during initialization:

```go
validate.RegisterValidation("strong_password", strongPassword)
validate.RegisterValidation("custom_url", isURL)
validate.RegisterValidation("date", isDate)
validate.RegisterValidation("time", isTime)
validate.RegisterValidation("ip_address", isIPAddress)
```

## Handler Integration

### Basic Usage Pattern

```go
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var userDTO model.RegisterRequest

    // Validate request with field errors
    if fieldErrors := util.ValidateRequestWithFields(r, &userDTO); fieldErrors != nil {
        response.ValidationErrorWithFields(w, fieldErrors)
        return
    }

    // If validation passes, continue with business logic...
    user, err := h.authService.Register(r.Context(), &userDTO)
    // ...
}
```

### Under the Hood

1. `ValidateRequestWithFields` decodes the JSON body into your struct
2. It then validates the struct against its validation tags
3. If there are validation errors, it returns a map of `fieldName -> errorMessage`
4. `response.ValidationErrorWithFields` formats this map into a standardized JSON response

## Error Responses

### Field-Based Error Response Format

When validation fails, the API returns a response like this:

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "fields": {
      "email": "email is required",
      "password": "password must be at least 8 characters long",
      "firstName": "firstName is required",
      "lastName": "lastName is required"
    }
  }
}
```

### Response Structure

- `success`: Always `false` for error responses
- `error`: Contains error details
  - `code`: Error code, always `VALIDATION_ERROR` for validation errors
  - `message`: General error message
  - `fields`: Map of field name to error message

## Best Practices

### 1. Keep Validation Close to Your Models

Define validation tags directly on your DTOs:

```go
// Good: Model with built-in validation
type RegisterRequest struct {
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8,strong_password"`
    FirstName string `json:"first_name" validate:"required"`
    LastName  string `json:"last_name" validate:"required"`
}
```

### 2. Use Consistent Validation Patterns

All handlers should follow the same validation pattern using `ValidateRequestWithFields`:

```go
// Standard validation pattern
if fieldErrors := util.ValidateRequestWithFields(r, &dto); fieldErrors != nil {
    response.ValidationErrorWithFields(w, fieldErrors)
    return
}
```

This pattern provides field-specific error messages that are more helpful to API consumers than generic validation errors.

### 3. Group Related Validations

For complex validations, create custom validators:

```go
// Instead of multiple separate tags
Password string `json:"password" validate:"min=8,containsUpper,containsLower,containsNumber,containsSpecial"`

// Use a single custom validator
Password string `json:"password" validate:"strong_password"`
```

### 4. Add Documentation for Domain-Specific Validations

```go
// ProductRequest defines a product creation request
type ProductRequest struct {
    // SKU must follow pattern ABC-12345
    SKU string `json:"sku" validate:"required,sku_format"`

    // Price must be greater than 0 and have at most 2 decimal places
    Price float64 `json:"price" validate:"required,gt=0"`
}
```

## Examples

### Complete Registration Example

**Model:**

```go
type RegisterRequest struct {
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8,strong_password"`
    FirstName string `json:"first_name" validate:"required"`
    LastName  string `json:"last_name" validate:"required"`
    Age       int    `json:"age" validate:"required,min=18"`
    Country   string `json:"country" validate:"required"`
}
```

**Handler:**

```go
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var userDTO model.RegisterRequest

    // Validate request with field errors
    if fieldErrors := util.ValidateRequestWithFields(r, &userDTO); fieldErrors != nil {
        response.ValidationErrorWithFields(w, fieldErrors)
        return
    }

    user, err := h.authService.Register(r.Context(), &userDTO)
    if err != nil {
        if errors.Is(err, service.ErrEmailAlreadyExists) {
            response.Error(w, http.StatusConflict, "User with this email already exists", err.Error())
            return
        }
        response.Error(w, http.StatusInternalServerError, "Failed to register user", err.Error())
        return
    }

    response.JSON(w, http.StatusCreated, user)
}
```

### Invalid Request and Response

**Request:**

```json
{
  "email": "invalid-email",
  "password": "weak",
  "first_name": "",
  "age": 16
}
```

**Response:**

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "fields": {
      "email": "email must be a valid email address",
      "password": "password must be at least 8 characters long",
      "password": "password must contain at least 1 uppercase, 1 lowercase, 1 number, and 1 special character",
      "firstName": "firstName is required",
      "lastName": "lastName is required",
      "age": "age must be at least 18"
    }
  }
}
```

## Technical Details

### Understanding the Underscore `_` in Go

In the validator registration code, you'll notice lines like this:

```go
_ = validate.RegisterValidation("strong_password", strongPassword)
```

The underscore `_` here is a special identifier in Go that indicates we're intentionally discarding the return value. Here's what's happening:

1. `RegisterValidation` returns an error if the registration fails
2. The `_` tells Go that we're aware of this return value, but we've chosen not to handle it
3. Without the `_`, Go would give a compiler error because we're not using the returned error

This is a common pattern in Go when:

- You expect the operation to succeed in normal circumstances
- You're not planning to handle errors from this particular function
- The application can continue even if this specific operation fails

If you wanted to handle errors, you could replace the `_` with a variable name:

```go
if err := validate.RegisterValidation("strong_password", strongPassword); err != nil {
    // Handle the error, perhaps by logging it or panicking
    log.Fatalf("Failed to register validator: %v", err)
}
```

For our validator initialization, we're using the simpler approach because:

1. Registration failures are rare and would indicate a programming error
2. We're registering these validators during initialization, not in response to user input
3. If registration did fail, there's not much we could do to recover at runtime

This comprehensive validation system makes it easy to:

1. Define validation rules directly on your models
2. Validate requests in your handlers with minimal code
3. Return consistent, field-specific error responses
4. Create custom validators for domain-specific rules

For more validation options, refer to the [go-playground/validator documentation](https://pkg.go.dev/github.com/go-playground/validator/v10).
