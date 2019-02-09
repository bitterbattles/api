package http

// Request represents an HTTP request
type Request struct {
	QueryParams map[string]string `json:"queryStringParameters"`
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body"`
}
