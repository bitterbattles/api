package http

import "strings"

// Request represents an HTTP request
type Request struct {
	PathParams  map[string]string `json:"pathParameters"`
	QueryParams map[string]string `json:"queryStringParameters"`
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body"`
}

// GetHeader gets a header using case-insensitive compare
func (request *Request) GetHeader(header string) string {
	value := request.Headers[header]
	if value == "" {
		value = request.Headers[strings.ToLower(header)]
	}
	return value
}
