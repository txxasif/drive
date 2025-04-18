package util

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

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

	// Register function to get JSON tag as field name
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

// ValidateStructWithFields validates a struct and returns validation errors with field mappings
func ValidateStructWithFields(i interface{}) map[string]string {
	if err := validate.Struct(i); err == nil {
		return nil
	} else {
		// log.Println("Validation error:", err)

		fieldErrors := make(map[string]string)
		// log.Println("fieldErrors", err.(validator.ValidationErrors))
		// Collect validation errors
		if validatorErrs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validatorErrs {
				field := e.Field()
				log.Println("field", field)
				fieldErrors[field] = field + " " + formatValidationErrorMessage(e)
			}
		} else {
			fieldErrors["_general"] = err.Error()
		}

		return fieldErrors
	}
}

// ValidateRequestWithFields decodes the request body and returns field-based validation errors
func ValidateRequestWithFields(r *http.Request, dst interface{}) map[string]string {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return map[string]string{"_general": "Invalid request body: " + err.Error()}
	}
	return ValidateStructWithFields(dst)
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

// isIPAddress checks if a string is a valid IP address
func isIPAddress(fl validator.FieldLevel) bool {
	ip := fl.Field().String()
	return ipv4Regex.MatchString(ip) || ipv6Regex.MatchString(ip)
}
