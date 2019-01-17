package input

import (
	"fmt"
	"strings"
)

// StringRules describes rules for validating/sanitizing strings
type StringRules struct {
	ToLower      bool
	TrimSpace    bool
	Length       int
	MinLength    int
	MaxLength    int
	ValidValues  []string
	DefaultValue string
}

// SanitizeString performs validation and sanitization on a string value
func SanitizeString(value string, rules StringRules, errorCreator ErrorCreator) (string, error) {
	if rules.ToLower {
		value = strings.ToLower(value)
	}
	if rules.TrimSpace {
		value = strings.TrimSpace(value)
	}
	if rules.ValidValues != nil {
		if isValidValue(value, rules.ValidValues) {
			return value, nil
		}
		if rules.DefaultValue != "" {
			return rules.DefaultValue, nil
		}
		return "", errorCreator(fmt.Sprintf("Not one of the allowed values."))
	}
	length := len(value)
	if rules.Length > 0 && length != rules.Length {
		return "", errorCreator(fmt.Sprintf("Length must be exactly %d.", rules.Length))
	}
	if rules.MinLength > 0 && length < rules.MinLength {
		return "", errorCreator(fmt.Sprintf("Minimum length is %d.", rules.MinLength))
	}
	if rules.MaxLength > 0 && length > rules.MaxLength {
		return "", errorCreator(fmt.Sprintf("Maximum length is %d.", rules.MaxLength))
	}
	return value, nil
}

func isValidValue(value string, validValues []string) bool {
	for _, validValue := range validValues {
		if value == validValue {
			return true
		}
	}
	return false
}
