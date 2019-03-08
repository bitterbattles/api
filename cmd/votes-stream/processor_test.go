package main_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/bitterbattles/api/cmd/votes-stream"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/lambda/stream"
	. "github.com/bitterbattles/api/pkg/tests"
)

func TestProcessor(t *testing.T) {
	repository := mocks.NewRepository()
	addBattles(repository, 2)
	processor := NewProcessor(repository)
	eventJSON := `{"records":[{"dynamodb":{"NewImage":{"battleID":{"S":"id0"},"isVoteFor":{"BOOL":true}}}},{"dynamodb":{"NewImage":{"battleID":{"S":"id1"},"isVoteFor":{"BOOL":true}}}},{"dynamodb":{"NewImage":{"battleID":{"S":"id0"},"isVoteFor":{"BOOL":false}}}},{"dynamodb":{"NewImage":{"battleID":{"S":"id0"},"isVoteFor":{"BOOL":true}}}}]}`
	event := newEvent(eventJSON)
	err := processor.Process(event)
	AssertNil(t, err)
	verifyBattleVotes(t, repository, "id0", 2, 1)
	verifyBattleVotes(t, repository, "id1", 1, 0)
}

func addBattles(repository *mocks.Repository, count int) {
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

func verifyBattleVotes(t *testing.T, repository *mocks.Repository, id string, expectedVotesFor int, expectedVotesAgainst int) {
	battle, err := repository.GetByID(id)
	AssertNil(t, err)
	AssertNotNil(t, battle)
	AssertEquals(t, battle.VotesFor, expectedVotesFor)
	AssertEquals(t, battle.VotesAgainst, expectedVotesAgainst)
}
