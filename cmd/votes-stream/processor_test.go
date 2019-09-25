package main_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/bitterbattles/api/cmd/votes-stream"
	"github.com/bitterbattles/api/pkg/battles"
	battlesMocks "github.com/bitterbattles/api/pkg/battles/mocks"
	indexMocks "github.com/bitterbattles/api/pkg/index/mocks"
	"github.com/bitterbattles/api/pkg/lambda/stream"
	. "github.com/bitterbattles/api/pkg/tests"
)

func TestProcessor(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	indexer := battles.NewIndexer(indexRepository)
	battlesRepository := battlesMocks.NewRepository()
	addBattles(battlesRepository, 2)
	processor := NewProcessor(indexer, battlesRepository)
	eventJSON := `{"records":[{"dynamodb":{"NewImage":{"userId":{"S":"userId0"},"battleId":{"S":"id0"},"isVoteFor":{"BOOL":true}}}},{"dynamodb":{"NewImage":{"userId":{"S":"userId0"},"battleId":{"S":"id1"},"isVoteFor":{"BOOL":true}}}},{"dynamodb":{"NewImage":{"userId":{"S":"userId1"},"battleId":{"S":"id0"},"isVoteFor":{"BOOL":false}}}},{"dynamodb":{"NewImage":{"userId":{"S":"userId2"},"battleId":{"S":"id0"},"isVoteFor":{"BOOL":true}}}}]}`
	event := newEvent(eventJSON)
	err := processor.Process(event)
	AssertNil(t, err)
	verifyBattleVotes(t, battlesRepository, "id0", 2, 1)
	verifyBattleVotes(t, battlesRepository, "id1", 1, 0)
}

func addBattles(repository *battlesMocks.Repository, count int) {
	for i := 0; i < count; i++ {
		battle := battles.Battle{
			ID:           fmt.Sprintf("id%d", i),
			VotesFor:     0,
			VotesAgainst: 0,
		}
		repository.Add(&battle)
	}
}

func newEvent(eventJSON string) *stream.Event {
	event := &stream.Event{}
	json.Unmarshal([]byte(eventJSON), event)
	return event
}

func verifyBattleVotes(t *testing.T, repository *battlesMocks.Repository, id string, expectedVotesFor int, expectedVotesAgainst int) {
	battle, err := repository.GetByID(id)
	AssertNil(t, err)
	AssertNotNil(t, battle)
	AssertEquals(t, battle.VotesFor, expectedVotesFor)
	AssertEquals(t, battle.VotesAgainst, expectedVotesAgainst)
}
