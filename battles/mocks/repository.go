package mocks

import (
	"errors"

	"github.com/bitterbattles/api/battles"
)

// Repository is a mocked implementation of battles.RepositoryInterface
type Repository struct {
	data      map[string]*battles.Battle
	lastAdded *battles.Battle
}

// NewRepository creates a new Battles repository instance
func NewRepository() *Repository {
	return &Repository{make(map[string]*battles.Battle), nil}
}

// Add is used to insert a new Battle document
func (repository *Repository) Add(battle battles.Battle) error {
	repository.data[battle.ID] = &battle
	repository.lastAdded = &battle
	return nil
}

// GetByID is used to get a Battle by ID
func (repository *Repository) GetByID(id string) (*battles.Battle, error) {
	_, ok := repository.data[id]
	if !ok {
		return nil, nil
	}
	return repository.data[id], nil
}

// IncrementVotes increments the votes for a given Battle ID
func (repository *Repository) IncrementVotes(id string, deltaVotesFor int, deltaVotesAgainst int) error {
	_, ok := repository.data[id]
	if !ok {
		return errors.New("battle not found")
	}
	battle := repository.data[id]
	battle.VotesFor += deltaVotesFor
	battle.VotesAgainst += deltaVotesAgainst
	repository.data[id] = battle
	return nil
}

// GetLastAdded gets the most recently added Battle
func (repository *Repository) GetLastAdded() *battles.Battle {
	return repository.lastAdded
}

// Reset removes all data from the repository
func (repository *Repository) Reset() {
	repository.data = make(map[string]*battles.Battle)
	repository.lastAdded = nil
}
