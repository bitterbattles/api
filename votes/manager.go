package votes

import (
	"time"

	"github.com/bitterbattles/api/core/input"
	"github.com/bitterbattles/api/votes/errors"
)

// ManagerInterface defines an interface for a Vote manager
type ManagerInterface interface {
	Create(string, bool) error
}

// Manager is used to perform business logic related to Votes
type Manager struct {
	repository RepositoryInterface
}

// NewManager creates a new Manager instance
func NewManager(repository RepositoryInterface) *Manager {
	return &Manager{repository}
}

// Create creates a new Vote
func (manager *Manager) Create(battleID string, isVoteFor bool) error {
	battleID, err := manager.sanitizeBattleID(battleID)
	if err != nil {
		return err
	}
	vote := Vote{
		UserID:    "bgttr132fopt0uo06vlg",
		BattleID:  battleID,
		IsVoteFor: isVoteFor,
		CreatedOn: time.Now().Unix(),
	}
	return manager.repository.Add(vote)
}

func (manager *Manager) sanitizeBattleID(battleID string) (string, error) {
	rules := input.StringRules{
		TrimSpace: true,
		Length:    battleIDLength,
	}
	errorCreator := func(message string) error {
		return errors.NewInvalidBattleIDError("Invalid Battle ID: " + message)
	}
	return input.SanitizeString(battleID, rules, errorCreator)
}
