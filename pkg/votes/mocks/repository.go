package mocks

import (
	"github.com/bitterbattles/api/pkg/votes"
)

// Repository is a mocked implementation of votes.RepositoryInterface
type Repository struct {
	lastAdded *votes.Vote
}

// NewRepository creates a new Votes repository instance
func NewRepository() *Repository {
	return &Repository{nil}
}

// Add is used to insert a new Vote document
func (repository *Repository) Add(vote votes.Vote) error {
	repository.lastAdded = &vote
	return nil
}

// GetLastAdded gets the most recently added Vote
func (repository *Repository) GetLastAdded() *votes.Vote {
	return repository.lastAdded
}

// Reset removes all data from the repository
func (repository *Repository) Reset() {
	repository.lastAdded = nil
}
