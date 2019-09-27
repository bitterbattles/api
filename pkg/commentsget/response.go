package commentsget

// Response a comment in a GET response
type Response struct {
	ID        string `json:"id"`
	BattleID  string `json:"battleId"`
	CreatedOn int64  `json:"createdOn"`
	Username  string `json:"username"`
	Comment   string `json:"comment"`
}
