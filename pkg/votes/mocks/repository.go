package mocks

import (
	"github.com/bitterbattles/api/pkg/votes"
)

// Repository is a mocked implementation of votes.RepositoryInterface
type Repository struct {
	data      map[string]*votes.Vote
	lastAdded *votes.Vote
}

// NewRepository creates a new Votes repository instance
func NewRepository() *Repository {
	return &Repository{make(map[string]*votes.Vote), nil}
}

// Add is used to insert a new Vote document
func (repository *Repository) Add(vote *votes.Vote) error {
	key := vote.UserID + "-" + vote.BattleID
	repository.data[key] = vote
	repository.lastAdded = vote
	return nil
}

// GetByUserAndBattleIDs is used to get a Vote by user ID and battle ID
func (repository *Repository) GetByUserAndBattleIDs(userID string, battleID string) (*votes.Vote, error) {
	key := userID + "-" + battleID
	return repository.data[key], nil
}

// GetLastAdded gets the most recently added Vote
func (repository *Repository) GetLastAdded() *votes.Vote {
	return repository.lastAdded
}

// Reset removes all data from the repository
func (repository *Repository) Reset() {
	repository.data = make(map[string]*votes.Vote)
	repository.lastAdded = nil
}
