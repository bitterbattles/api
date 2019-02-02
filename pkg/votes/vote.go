package votes

// Vote model
type Vote struct {
	UserID    string `json:"userId"`
	BattleID  string `json:"battleId"`
	IsVoteFor bool   `json:"isVoteFor"`
	CreatedOn int64  `json:"createdOn"`
}
