package main_test

import (
	"encoding/json"
	"testing"

	. "github.com/bitterbattles/api/cmd/battles-stream"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/lambda/stream"
	. "github.com/bitterbattles/api/pkg/tests"
)

func TestProcessorNewBattle(t *testing.T) {
	index := mocks.NewIndex()
	processor := NewProcessor(index, battles.NewScorer())
	eventJSON := `{"records":[{"dynamodb":{"NewImage":{"id":{"S":"id0"},"userId":{"S":"userId0"},"title":{"S":"title0"},"description":{"S":"description0"},"votesFor":{"N":"1"},"votesAgainst":{"N":"2"},"comments":{"N":"3"},"createdOn":{"N":"86400"},"state":{"N":"1"}}}}]}`
	event := newEvent(eventJSON)
	err := processor.Process(event)
	AssertNil(t, err)
	expectedBattle := &battles.Battle{
		ID:           "id0",
		UserID:       "userId0",
		Title:        "title0",
		Description:  "description0",
		VotesFor:     1,
		VotesAgainst: 2,
		Comments:     3,
		CreatedOn:    86400,
		State:        1,
	}
	verifyIndex(t, index, expectedBattle, 7, 6)
}

func TestProcessorUpdatedBattle(t *testing.T) {
	index := mocks.NewIndex()
	processor := NewProcessor(index, battles.NewScorer())
	eventJSON := `{"records":[{"dynamodb":{"OldImage":{"id":{"S":"id0"},"userId":{"S":"userId0"},"title":{"S":"title0"},"description":{"S":"description0"},"votesFor":{"N":"1"},"votesAgainst":{"N":"2"},"comments":{"N":"3"},"createdOn":{"N":"86400"},"state":{"N":"1"}},"NewImage":{"id":{"S":"id0"},"userId":{"S":"userId0"},"title":{"S":"title0"},"description":{"S":"description0"},"votesFor":{"N":"4"},"votesAgainst":{"N":"5"},"comments":{"N":"6"},"createdOn":{"N":"86400"},"state":{"N":"1"}}}}]}`
	event := newEvent(eventJSON)
	err := processor.Process(event)
	AssertNil(t, err)
	expectedBattle := &battles.Battle{
		ID:           "id0",
		UserID:       "userId0",
		Title:        "title0",
		Description:  "description0",
		VotesFor:     4,
		VotesAgainst: 5,
		Comments:     6,
		CreatedOn:    86400,
		State:        1,
	}
	verifyIndex(t, index, expectedBattle, 16, 15)
}

func TestProcessorDeletedBattle(t *testing.T) {
	index := mocks.NewIndex()
	processor := NewProcessor(index, battles.NewScorer())
	eventJSON := `{"records":[{"dynamodb":{"OldImage":{"id":{"S":"id0"},"userId":{"S":"userId0"},"title":{"S":"title0"},"description":{"S":"description0"},"votesFor":{"N":"1"},"votesAgainst":{"N":"2"},"comments":{"N":"3"},"createdOn":{"N":"123"},"state":{"N":"1"}},"NewImage":{"id":{"S":"id0"},"userId":{"S":"userId0"},"title":{"S":"title0"},"description":{"S":"description0"},"votesFor":{"N":"1"},"votesAgainst":{"N":"2},"comments":{"N":"3"},"createdOn":{"N":"123"},"state":{"N":"2"}}}}]}`
	event := newEvent(eventJSON)
	err := processor.Process(event)
	AssertNil(t, err)
	verifyIndex(t, index, nil, 0, 0)
}

func newEvent(eventJSON string) *stream.Event {
	event := &stream.Event{}
	json.Unmarshal([]byte(eventJSON), event)
	return event
}
func verifyIndex(t *testing.T, index *mocks.Index, expectedBattle *battles.Battle, expectedPopularity int64, expectedControversy int64) {
	battle := index.GetLastBattleAdded()
	popularity := index.GetLastPopularityAdded()
	controversy := index.GetLastControversyAdded()
	if expectedBattle == nil {
		AssertNil(t, battle)
		AssertEquals(t, popularity, expectedPopularity)
		AssertEquals(t, controversy, expectedControversy)
	} else {
		AssertNotNil(t, battle)
		AssertEquals(t, battle.ID, expectedBattle.ID)
		AssertEquals(t, battle.UserID, expectedBattle.UserID)
		AssertEquals(t, battle.Title, expectedBattle.Title)
		AssertEquals(t, battle.Description, expectedBattle.Description)
		AssertEquals(t, battle.VotesFor, expectedBattle.VotesFor)
		AssertEquals(t, battle.VotesAgainst, expectedBattle.VotesAgainst)
		AssertEquals(t, battle.Comments, expectedBattle.Comments)
		AssertEquals(t, battle.CreatedOn, expectedBattle.CreatedOn)
		AssertEquals(t, popularity, expectedPopularity)
		AssertEquals(t, controversy, expectedControversy)
	}
}
