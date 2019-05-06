package battlesget

import (
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/users"
)

const (
	defaultUsername = "[Deleted]"
)

// Response a Battle in a GET response
type Response struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	CanVote      bool   `json:"canVote"`
	VotesFor     int    `json:"votesFor"`
	VotesAgainst int    `json:"votesAgainst"`
	CreatedOn    int64  `json:"createdOn"`
}

// ToResponse converts a Battle to a GET response
func ToResponse(battle *battles.Battle, user *users.User, canVote bool) *Response {
	username := defaultUsername
	if user != nil && user.State == users.Active {
		username = user.DisplayUsername
	}
	return &Response{
		ID:           battle.ID,
		Username:     username,
		Title:        battle.Title,
		Description:  battle.Description,
		CanVote:      canVote,
		VotesFor:     battle.VotesFor,
		VotesAgainst: battle.VotesAgainst,
		CreatedOn:    battle.CreatedOn,
	}
}
