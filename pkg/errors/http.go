package errors

import "github.com/bitterbattles/api/pkg/http"

// HTTPError represents an error appropriate for HTTP responses
type HTTPError struct {
	statusCode   int
	errorCode    int
	errorMessage string
}

// StatusCode gets the status code
func (error HTTPError) StatusCode() int {
	return error.statusCode
}

// ErrorCode gets the status code
func (error HTTPError) ErrorCode() int {
	return error.errorCode
}

// Error gets the error message (complies with error interface)
func (error HTTPError) Error() string {
	return error.errorMessage
}

// NewBadRequestError creates a new HTTPError instance representing a Bad Request
func NewBadRequestError(message string) *HTTPError {
	return NewBadRequestErrorWithCode(http.BadRequest, message)
}

// NewBadRequestErrorWithCode creates a new HTTPError instance representing a Bad Request with a specific error code
func NewBadRequestErrorWithCode(code int, message string) *HTTPError {
	return newHTTPError(http.BadRequest, code, message)
}

// NewForbiddenError creates a new HTTPError instance representing a Forbidden error
func NewForbiddenError(message string) *HTTPError {
	return newHTTPError(http.Forbidden, http.Forbidden, message)
}

// NewNotFoundError creates a new HTTPError instance representing a Not Found error
func NewNotFoundError(message string) *HTTPError {
	return newHTTPError(http.NotFound, http.NotFound, message)
}

func newHTTPError(statusCode int, errorCode int, errorMessage string) *HTTPError {
	return &HTTPError{
		statusCode:   statusCode,
		errorCode:    errorCode,
		errorMessage: errorMessage,
	}
}
