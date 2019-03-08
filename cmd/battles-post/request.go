package main

// Request represents a new Battle request
type Request struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
