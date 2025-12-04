package utils

import (
	"regexp"
	"strings"
)

// ValidateMSISDN validates Bangladesh mobile number format (must start with 01)
func ValidateMSISDN(msisdn string) bool {
	// Remove any spaces and special characters
	cleaned := strings.ReplaceAll(msisdn, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	// Bangladesh mobile number pattern - must start with 01 and have 11 digits total
	// Valid formats: 015xxxxxxxx, 016xxxxxxxx, 017xxxxxxxx, 018xxxxxxxx, 019xxxxxxxx
	pattern := `^01[3-9]\d{8}$`

	matched, _ := regexp.MatchString(pattern, cleaned)
	return matched
}

// NormalizeMSISDN normalizes MSISDN to standard BD format (01xxxxxxxxx)
func NormalizeMSISDN(msisdn string) string {
	// Remove any spaces and special characters
	cleaned := strings.ReplaceAll(msisdn, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	// If already in correct format (starts with 01)
	if matched, _ := regexp.MatchString(`^01[3-9]\d{8}$`, cleaned); matched {
		return cleaned
	}

	// If starts with +8801, remove the prefix
	if matched, _ := regexp.MatchString(`^\+8801[3-9]\d{8}$`, cleaned); matched {
		return cleaned[3:] // Remove "+88" prefix
	}

	// If starts with 8801, remove the prefix
	if matched, _ := regexp.MatchString(`^8801[3-9]\d{8}$`, cleaned); matched {
		return cleaned[2:] // Remove "88" prefix
	}

	return cleaned // Return as is (will be validated separately)
}

// ValidatePassword validates password strength
func ValidatePassword(password string) bool {
	// Minimum 6 characters
	if len(password) < 6 {
		return false
	}

	// Add more complex validation if needed
	// For now, just check length
	return true
}

// ValidateName validates user name
func ValidateName(name string) bool {
	name = strings.TrimSpace(name)

	// Must be between 2 and 100 characters
	if len(name) < 2 || len(name) > 100 {
		return false
	}

	// Should not contain only numbers or special characters
	matched, _ := regexp.MatchString(`^[a-zA-Z\s.'-]+$`, name)
	return matched
}
