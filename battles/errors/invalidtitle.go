package errors

// InvalidTitleError represents an error occurring due to an invalid title
type InvalidTitleError struct {
	message string
}

// NewInvalidTitleError creates a new InvalidTitleError instance
func NewInvalidTitleError(message string) InvalidTitleError {
	return InvalidTitleError{message}
}

// Error gets the error message
func (error InvalidTitleError) Error() string {
	return error.message
}
