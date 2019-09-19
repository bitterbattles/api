package battlesget

import (
	"log"

	"github.com/bitterbattles/api/pkg/battles"
)

// CreateResponses creates a list of GET battles responses
func CreateResponses(userID string, battleIDs []string, repository battles.RepositoryInterface, getCanVote func(string, *battles.Battle) (bool, error)) ([]*Response, error) {
	responses := make([]*Response, 0, len(battleIDs))
	for _, battleID := range battleIDs {
		battle, err := repository.GetByID(battleID)
		if err != nil {
			return nil, err
		}
		if battle == nil {
			log.Println("Failed to find battle with ID", battleID, "when fetching battles.")
			continue
		}
		if battle.State == battles.Deleted {
			continue
		}
		canVote, err := getCanVote(userID, battle)
		if err != nil {
			return nil, err
		}
		response := &Response{
			ID:           battle.ID,
			Username:     battle.Username,
			Title:        battle.Title,
			Description:  battle.Description,
			CanVote:      canVote,
			VotesFor:     battle.VotesFor,
			VotesAgainst: battle.VotesAgainst,
			CreatedOn:    battle.CreatedOn,
		}
		responses = append(responses, response)
	}
	return responses, nil
}
