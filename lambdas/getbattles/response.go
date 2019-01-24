package getbattles

// Response represents a response
type Response struct {
	ID           string `json:"id"`
	UserID       string `json:"-"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	VotesFor     int    `json:"votesFor"`
	VotesAgainst int    `json:"votesAgainst"`
	CreatedOn    int64  `json:"createdOn"`
}
