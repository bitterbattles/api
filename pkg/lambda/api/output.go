package api

import "github.com/bitterbattles/api/pkg/http"

// Output represents the result returned from a processor
type Output struct {
	StatusCode   int
	ResponseBody interface{}
}

// NewOutput creats a new Output instance with status code OK
func NewOutput(responseBody interface{}) *Output {
	return NewOutputWithStatus(http.OK, responseBody)
}

// NewOutputWithStatus creats a new Output instance with specified status code
func NewOutputWithStatus(statusCode int, responseBody interface{}) *Output {
	return &Output{
		StatusCode:   statusCode,
		ResponseBody: responseBody,
	}
}
