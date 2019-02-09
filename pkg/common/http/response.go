package http

import (
	"encoding/json"
)

// Response represents an HTTP response
type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// NewResponse creates a new Response instance with status code OK
func NewResponse(body interface{}, headers map[string]string) (*Response, error) {
	return NewResponseWithStatus(body, headers, OK)
}

// NewResponseWithStatus creates a new Response instance including a custom status code
func NewResponseWithStatus(body interface{}, headers map[string]string, statusCode int) (*Response, error) {
	response := &Response{
		StatusCode: statusCode,
		Headers:    headers,
	}
	if body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		response.Body = string(bodyJSON)
	}
	return response, nil
}
