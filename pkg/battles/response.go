package battles

import (
	"github.com/bitterbattles/api/pkg/users"
)

const (
	defaultUsername = "[Deleted]"
)

// Response represents an element in the GET response
type Response struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	HasVoted     bool   `json:"hasVoted"`
	VotesFor     int    `json:"votesFor"`
	VotesAgainst int    `json:"votesAgainst"`
	CreatedOn    int64  `json:"createdOn"`
}

// ToGetResponse converts a Battle to a GET response
func ToGetResponse(battle *Battle, user *users.User, hasVoted bool) *Response {
	username := defaultUsername
	if user != nil && user.State == users.Active {
		username = user.DisplayUsername
	}
	return &Response{
		ID:           battle.ID,
		Username:     username,
		Title:        battle.Title,
		Description:  battle.Description,
		HasVoted:     hasVoted,
		VotesFor:     battle.VotesFor,
		VotesAgainst: battle.VotesAgainst,
		CreatedOn:    battle.CreatedOn,
	}
}
