package errors

// HTTPError represents an error appropriate for HTTP responses
type HTTPError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// Error gets the error message
func (error HTTPError) Error() string {
	return error.Message
}

// NewBadRequestError creates a new HTTPError instance representing a Bad Request
func NewBadRequestError(message string) HTTPError {
	return HTTPError{BadRequestCode, message}
}

// NewInternalServerError creates a new HTTPError instance representing an Internal Server Error
func NewInternalServerError() HTTPError {
	return HTTPError{InternalServerErrorCode, "Something unexpected happened. Please try again later."}
}
