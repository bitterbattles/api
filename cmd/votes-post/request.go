package main

// Request represents a request body
type Request struct {
	BattleID  string `json:"battleId"`
	IsVoteFor bool   `json:"isVoteFor"`
}
