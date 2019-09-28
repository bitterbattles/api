package battlesget

// Response a Battle in a GET response
type Response struct {
	ID           string `json:"id"`
	CreatedOn    int64  `json:"createdOn"`
	Username     string `json:"username"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	CanVote      bool   `json:"canVote"`
	VotesFor     int    `json:"votesFor"`
	VotesAgainst int    `json:"votesAgainst"`
	Comments     int    `json:"comments"`
	Verdict      int    `json:"verdict"`
}
