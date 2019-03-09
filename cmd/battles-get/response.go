package main

// Response represents an element in the response results
type Response struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	HasVoted     bool   `json:"hasVoted"`
	VotesFor     int    `json:"votesFor"`
	VotesAgainst int    `json:"votesAgainst"`
	CreatedOn    int64  `json:"createdOn"`
}
