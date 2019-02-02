package votes

// Vote model
type Vote struct {
	UserID    string `json:"userID"`
	BattleID  string `json:"battleID"`
	IsVoteFor bool   `json:"isVoteFor"`
	CreatedOn int64  `json:"createdOn"`
}
