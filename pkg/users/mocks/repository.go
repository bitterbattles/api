package mocks

import (
	"strings"

	"github.com/bitterbattles/api/pkg/users"
)

// Repository is a mocked implementation of users.RepositoryInterface
type Repository struct {
	idData       map[string]*users.User
	usernameData map[string]*users.User
	lastAdded    *users.User
}

// NewRepository creates a new User repository instance
func NewRepository() *Repository {
	return &Repository{
		idData:       make(map[string]*users.User),
		usernameData: make(map[string]*users.User),
		lastAdded:    nil,
	}
}

// Add is used to insert a new User document
func (repository *Repository) Add(user *users.User) error {
	repository.idData[user.ID] = user
	username := strings.ToLower(user.Username)
	repository.usernameData[username] = user
	repository.lastAdded = user
	return nil
}

// GetByID is used to get a User by ID
func (repository *Repository) GetByID(id string) (*users.User, error) {
	return repository.idData[id], nil
}

// GetByUsername looks up a User with the specified username
func (repository *Repository) GetByUsername(username string) (*users.User, error) {
	username = strings.ToLower(username)
	return repository.usernameData[username], nil
}

// GetLastAdded gets the most recently added User
func (repository *Repository) GetLastAdded() *users.User {
	return repository.lastAdded
}

// Reset removes all data from the repository
func (repository *Repository) Reset() {
	repository.idData = make(map[string]*users.User)
	repository.usernameData = make(map[string]*users.User)
	repository.lastAdded = nil
}
