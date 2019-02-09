package errors

import "github.com/bitterbattles/api/pkg/common/http"

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
	return HTTPError{http.BadRequest, message}
}
