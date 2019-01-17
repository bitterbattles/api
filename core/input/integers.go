package input

import "fmt"

// IntegerRules describes rules for validating/sanitizing integers
type IntegerRules struct {
	MinValue           int
	EnforceMinValue    bool
	DefaultMinValue    int
	UseDefaultMinValue bool
	MaxValue           int
	EnforceMaxValue    bool
	DefaultMaxValue    int
	UseDefaultMaxValue bool
}

// SanitizeInteger validates and sanitizes integer values
func SanitizeInteger(value int, rules IntegerRules, errorCreator ErrorCreator) (int, error) {
	if rules.EnforceMinValue && value < rules.MinValue {
		if rules.UseDefaultMinValue {
			return rules.DefaultMinValue, nil
		}
		return 0, errorCreator(fmt.Sprintf("Minimum value is %d.", rules.MinValue))
	}
	if rules.EnforceMaxValue && value > rules.MaxValue {
		if rules.UseDefaultMaxValue {
			return rules.DefaultMaxValue, nil
		}
		return 0, errorCreator(fmt.Sprintf("Maximum value is %d.", rules.MaxValue))
	}
	return value, nil
}
