package main

// Request represents a request body
type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
