package comments

// Comment model
type Comment struct {
	ID        string `json:"id"`
	BattleID  string `json:"battleId"`
	UserID    string `json:"userId"`
	CreatedOn int64  `json:"createdOn"`
	Username  string `json:"username"`
	Comment   string `json:"comment"`
	State     int    `json:"state"`
}
