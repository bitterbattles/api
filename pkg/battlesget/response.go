package battlesget

// Response a Battle in a GET response
type Response struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	CanVote      bool   `json:"canVote"`
	VotesFor     int    `json:"votesFor"`
	VotesAgainst int    `json:"votesAgainst"`
	CreatedOn    int64  `json:"createdOn"`
}
