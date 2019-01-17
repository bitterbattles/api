package errors

// InvalidDescriptionError represents an error occurring due to an invalid title
type InvalidDescriptionError struct {
	message string
}

// NewInvalidDescriptionError creates a new InvalidDescriptionError instance
func NewInvalidDescriptionError(message string) InvalidDescriptionError {
	return InvalidDescriptionError{message}
}

// Error gets the error message
func (error InvalidDescriptionError) Error() string {
	return error.message
}
