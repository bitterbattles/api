package mocks

import (
	"strings"

	"github.com/bitterbattles/api/pkg/users"
)

// Repository is a mocked implementation of users.RepositoryInterface
type Repository struct {
	data      map[string]*users.User
	lastAdded *users.User
}

// NewRepository creates a new User repository instance
func NewRepository() *Repository {
	return &Repository{make(map[string]*users.User), nil}
}

// Add is used to insert a new User document
func (repository *Repository) Add(user *users.User) error {
	username := strings.ToLower(user.Username)
	repository.data[username] = user
	repository.lastAdded = user
	return nil
}

// GetByUsername looks up a User with the specified username
func (repository *Repository) GetByUsername(username string) (*users.User, error) {
	username = strings.ToLower(username)
	_, ok := repository.data[username]
	if !ok {
		return nil, nil
	}
	return repository.data[username], nil
}

// GetLastAdded gets the most recently added User
func (repository *Repository) GetLastAdded() *users.User {
	return repository.lastAdded
}

// Reset removes all data from the repository
func (repository *Repository) Reset() {
	repository.data = make(map[string]*users.User)
	repository.lastAdded = nil
}
