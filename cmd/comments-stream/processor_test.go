package main_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/bitterbattles/api/cmd/comments-stream"
	"github.com/bitterbattles/api/pkg/battles"
	battlesMocks "github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/comments"
	indexMocks "github.com/bitterbattles/api/pkg/index/mocks"
	"github.com/bitterbattles/api/pkg/lambda/stream"
	. "github.com/bitterbattles/api/pkg/tests"
)

const battleKeyPattern = "commentIds:forBattle:%s"
const authorKeyPattern = "commentIds:forAuthor:%s"

func TestProcessorNewComments(t *testing.T) {
	indexRepostiory := indexMocks.NewRepository()
	indexer := comments.NewIndexer(indexRepostiory)
	battlesRepository := battlesMocks.NewRepository()
	battle := &battles.Battle{
		ID:       "battleId0",
		Comments: 0,
	}
	battlesRepository.Add(battle)
	processor := NewProcessor(indexer, battlesRepository)
	eventJSON := `{"records":[{"dynamodb":{"NewImage":{"id":{"S":"id0"},"battleId":{"S":"battleId0"},"userId":{"S":"userId0"},"createdOn":{"N":"123"}}}},{"dynamodb":{"NewImage":{"id":{"S":"id1"},"battleId":{"S":"battleId0"},"userId":{"S":"userId1"},"createdOn":{"N":"456"}}}}]}`
	event := newEvent(eventJSON)
	err := processor.Process(event)
	AssertNil(t, err)
	verifyIndexes(t, indexRepostiory, "id0", "battleId0", "userId0", 123)
	verifyIndexes(t, indexRepostiory, "id1", "battleId0", "userId1", 456)
	verifyBattleComments(t, battlesRepository, "battleId0", 2)
}

func newEvent(eventJSON string) *stream.Event {
	event := &stream.Event{}
	json.Unmarshal([]byte(eventJSON), event)
	return event
}

func verifyIndexes(t *testing.T, repository *indexMocks.Repository, commentID string, battleID string, userID string, expectedScore float64) {
	scoreByBattle := repository.GetScore(fmt.Sprintf(battleKeyPattern, battleID), commentID)
	scoreByAuthor := repository.GetScore(fmt.Sprintf(authorKeyPattern, userID), commentID)
	AssertEquals(t, scoreByBattle, expectedScore)
	AssertEquals(t, scoreByAuthor, expectedScore)
}

func verifyBattleComments(t *testing.T, repository *battlesMocks.Repository, id string, expectedComments int) {
	battle, err := repository.GetByID(id)
	AssertNil(t, err)
	AssertNotNil(t, battle)
	AssertEquals(t, battle.Comments, expectedComments)
}
