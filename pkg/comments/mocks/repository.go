package mocks

import (
	"github.com/bitterbattles/api/pkg/comments"
)

// Repository is a mocked implementation of comments.RepositoryInterface
type Repository struct {
	data      map[string]*comments.Comment
	lastAdded *comments.Comment
}

// NewRepository creates a new Comment repository instance
func NewRepository() *Repository {
	return &Repository{make(map[string]*comments.Comment), nil}
}

// Add is used to insert a new Comment document
func (repository *Repository) Add(comment *comments.Comment) error {
	repository.data[comment.ID] = comment
	repository.lastAdded = comment
	return nil
}

// DeleteByID deletes a Comment by ID
func (repository *Repository) DeleteByID(id string) error {
	repository.data[id] = nil
	return nil
}

// GetByID is used to get a Comment by ID
func (repository *Repository) GetByID(id string) (*comments.Comment, error) {
	return repository.data[id], nil
}

// UpdateUsername updates the specified user's username
func (repository *Repository) UpdateUsername(userID string, username string) error {
	for id, comment := range repository.data {
		if comment.UserID == userID {
			comment.Username = "[Deleted]"
			repository.data[id] = comment
		}
	}
	return nil
}

// GetLastAdded gets the most recently added Comment
func (repository *Repository) GetLastAdded() *comments.Comment {
	return repository.lastAdded
}

// Reset removes all data from the repository
func (repository *Repository) Reset() {
	repository.data = make(map[string]*comments.Comment)
	repository.lastAdded = nil
}
