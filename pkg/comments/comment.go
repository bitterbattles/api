package comments

// Comment model
type Comment struct {
	ID        string `json:"id"`
	BattleID  string `json:"battleId"`
	UserID    string `json:"userId"`
	Comment   string `json:"comment"`
	CreatedOn int64  `json:"createdOn"`
	State     int    `json:"state"`
}
