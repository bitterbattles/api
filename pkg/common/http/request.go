package http

// Request represents an HTTP request
type Request struct {
	PathParams  map[string]string `json:"pathParameters"`
	QueryParams map[string]string `json:"queryStringParameters"`
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body"`
}
