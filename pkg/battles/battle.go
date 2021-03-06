package battles

// Battle model
type Battle struct {
	ID           string `json:"id"`
	UserID       string `json:"userId"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	VotesFor     int    `json:"votesFor"`
	VotesAgainst int    `json:"votesAgainst"`
	Comments     int    `json:"comments"`
	CreatedOn    int64  `json:"createdOn"`
	State        int    `json:"state"`
}
