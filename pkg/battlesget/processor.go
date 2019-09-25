package battlesget

import (
	"log"
	"math"

	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/time"
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
			CreatedOn:    battle.CreatedOn,
			Username:     battle.Username,
			Title:        battle.Title,
			Description:  battle.Description,
			CanVote:      canVote,
			VotesFor:     battle.VotesFor,
			VotesAgainst: battle.VotesAgainst,
			Verdict:      determineVerdict(battle.CreatedOn, float64(battle.VotesFor), float64(battle.VotesAgainst)),
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func determineVerdict(createdOn int64, votesFor float64, votesAgainst float64) int {
	daysOld := (time.NowUnix() - createdOn) / 86400
	if daysOld < 1 {
		return None
	}
	totalVotes := votesFor + votesAgainst
	deltaVotes := math.Abs(votesFor - votesAgainst)
	if totalVotes == 0 || deltaVotes/totalVotes <= 0.05 {
		return NoDecision
	}
	if votesAgainst > votesFor {
		return Against
	}
	return For
}
