package main

// Response represents an element in the response results
type Response struct {
	ID           string `json:"id"`
	UserID       string `json:"-"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	VotesFor     int    `json:"votesFor"`
	VotesAgainst int    `json:"votesAgainst"`
	CreatedOn    int64  `json:"createdOn"`
	State        int    `json:"-"`
}
