package input

import (
	"fmt"
	"unicode"
)

// PasswordRules describes rules for validating/sanitizing passwords
type PasswordRules struct {
	MinLength  int
	MaxLength  int
	MinUpper   int
	MinNumbers int
	MinSymbols int
}

// SanitizePassword performs validation and sanitization on a string value
func SanitizePassword(value string, rules PasswordRules, errorCreator ErrorCreator) (string, error) {
	if value == "" {
		return "", errorCreator(fmt.Sprintf("Password is required but missing."))
	}
	length := len(value)
	if rules.MinLength > 0 && length < rules.MinLength {
		return "", errorCreator(fmt.Sprintf("Minimum length is %d.", rules.MinLength))
	}
	if rules.MaxLength > 0 && length > rules.MaxLength {
		return "", errorCreator(fmt.Sprintf("Maximum length is %d.", rules.MaxLength))
	}
	var numUpper, numNumbers, numSymbols int
	for _, c := range value {
		if unicode.IsUpper(c) {
			numUpper++
		} else if unicode.IsNumber(c) {
			numNumbers++
		} else if unicode.IsPunct(c) || unicode.IsSymbol(c) {
			numSymbols++
		}
	}
	if rules.MinUpper > 0 && numUpper < rules.MinUpper {
		return "", errorCreator(fmt.Sprintf("Minimum number of upper case is %d.", rules.MinUpper))
	}
	if rules.MinNumbers > 0 && numNumbers < rules.MinNumbers {
		return "", errorCreator(fmt.Sprintf("Minimum number of numbers is %d.", rules.MinNumbers))
	}
	if rules.MinSymbols > 0 && numSymbols < rules.MinSymbols {
		return "", errorCreator(fmt.Sprintf("Minimum number of symbols is %d.", rules.MinSymbols))
	}
	return value, nil
}
