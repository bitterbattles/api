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

// NewJSONResponse creates a new Response instance with status code OK
func NewJSONResponse(body interface{}, headers map[string]string) (*Response, error) {
	return NewJSONResponseWithStatus(body, headers, OK)
}

// NewJSONResponseWithStatus creates a new Response instance including a custom status code
func NewJSONResponseWithStatus(body interface{}, headers map[string]string, statusCode int) (*Response, error) {
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
		if response.Headers == nil {
			response.Headers = make(map[string]string)
		}
		response.Headers[ContentType] = ApplicationJSON
	}
	return response, nil
}
