package util

import (
	"encoding/json"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// hi
func init() {
	validate = validator.New()

	// Register function to get json tag as field name
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})

	// Register custom validators
	_ = validate.RegisterValidation("strong_password", strongPassword)
	_ = validate.RegisterValidation("custom_url", isURL)
	_ = validate.RegisterValidation("date", isDate)
	_ = validate.RegisterValidation("time", isTime)
	_ = validate.RegisterValidation("ip_address", isIPAddress)
}

// ValidateStruct validates a struct and returns validation errors
func ValidateStruct(i interface{}) []string {
	err := validate.Struct(i)
	if err == nil {
		return nil
	}

	var errorMessages []string

	// Collect validation errors
	validatorErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return []string{err.Error()}
	}

	for _, e := range validatorErrs {
		errorMessages = append(errorMessages, formatValidationError(e))
	}

	return errorMessages
}

// ValidateStructWithFields validates a struct and returns validation errors with field mappings
func ValidateStructWithFields(i interface{}) map[string]string {
	err := validate.Struct(i)
	if err == nil {
		return nil
	}

	fieldErrors := make(map[string]string)

	// Collect validation errors
	validatorErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		fieldErrors["_general"] = err.Error()
		return fieldErrors
	}

	for _, e := range validatorErrs {
		field := e.Field()
		// Include the field name in the error message for a complete message
		fieldErrors[field] = field + " " + formatValidationErrorMessage(e)
	}

	return fieldErrors
}

// ValidateRequest decodes the request body into the provided struct and validates it
func ValidateRequest(r *http.Request, dst interface{}) []string {
	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return []string{"Invalid request body: " + err.Error()}
	}

	// Validate the request body
	return ValidateStruct(dst)
}

// ValidateRequestWithFields decodes the request body and returns field-based validation errors
func ValidateRequestWithFields(r *http.Request, dst interface{}) map[string]string {
	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return map[string]string{"_general": "Invalid request body: " + err.Error()}
	}

	// Validate the request body
	return ValidateStructWithFields(dst)
}

// formatValidationError formats a validation error into a user-friendly message
func formatValidationError(err validator.FieldError) string {
	field := err.Field()
	message := formatValidationErrorMessage(err)
	return field + " " + message
}

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

// strongPassword validates if a password is strong
func strongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

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
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// isURL checks if a string is a valid URL
func isURL(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	r := regexp.MustCompile(`^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(?:/[a-zA-Z0-9\-\._~:/?#[\]@!$&'()*+,;=]*)?$`)
	return r.MatchString(url)
}

// isDate checks if a string is a valid date in YYYY-MM-DD format
func isDate(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	r := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !r.MatchString(date) {
		return false
	}
	// Additional logic to validate date values could be added here
	return true
}

// isTime checks if a string is a valid time in HH:MM:SS format
func isTime(fl validator.FieldLevel) bool {
	time := fl.Field().String()
	r := regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d):([0-5]\d)$`)
	return r.MatchString(time)
}

// isIPAddress checks if a string is a valid IP address
func isIPAddress(fl validator.FieldLevel) bool {
	ip := fl.Field().String()
	ipv4 := regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	ipv6 := regexp.MustCompile(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`)
	return ipv4.MatchString(ip) || ipv6.MatchString(ip)
}
