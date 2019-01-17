package errors

// InvalidBattleIDError represents an error occurring due to an invalid title
type InvalidBattleIDError struct {
	message string
}

// NewInvalidBattleIDError creates a new InvalidBattleIDError instance
func NewInvalidBattleIDError(message string) InvalidBattleIDError {
	return InvalidBattleIDError{message}
}

// Error gets the error message
func (error InvalidBattleIDError) Error() string {
	return error.message
}
