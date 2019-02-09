package main

// Request represents a request body
type Request struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
