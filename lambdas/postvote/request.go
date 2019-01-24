package postvote

// Request represents a request
type Request struct {
	BattleID  string `json:"battleId"`
	IsVoteFor bool   `json:"isVoteFor"`
}
