package main_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/bitterbattles/api/cmd/comments-stream"
	"github.com/bitterbattles/api/pkg/comments"
	"github.com/bitterbattles/api/pkg/index/mocks"
	"github.com/bitterbattles/api/pkg/lambda/stream"
	. "github.com/bitterbattles/api/pkg/tests"
)

const battleKeyPattern = "commentIds:forBattle:%s"
const authorKeyPattern = "commentIds:forAuthor:%s"

func TestProcessorNewComment(t *testing.T) {
	repository := mocks.NewRepository()
	indexer := comments.NewIndexer(repository)
	processor := NewProcessor(indexer)
	eventJSON := `{"records":[{"dynamodb":{"NewImage":{"id":{"S":"id0"},"battleId":{"S":"battleId0"},"userId":{"S":"userId0"},"createdOn":{"N":"123"}}}}]}`
	event := newEvent(eventJSON)
	err := processor.Process(event)
	AssertNil(t, err)
	verifyIndexes(t, repository, "id0", "battleId0", "userId0", 123)
}

func TestProcessorDeletedBattle(t *testing.T) {
	repository := mocks.NewRepository()
	setupIndexes(repository, "id0", "battleId0", "userId0", 123)
	indexer := comments.NewIndexer(repository)
	processor := NewProcessor(indexer)
	eventJSON := `{"records":[{"dynamodb":{"OldImage":{"id":{"S":"id0"},"battleId":{"S":"battleId0"},"userId":{"S":"userId0"},"createdOn":{"N":"123"},"state":{"N":"1"}},"NewImage":{"id":{"S":"id0"},"battleId":{"S":"battleId0"},"userId":{"S":"userId0"},"createdOn":{"N":"123"},"state":{"N":"2"}}}}]}`
	event := newEvent(eventJSON)
	err := processor.Process(event)
	AssertNil(t, err)
	verifyIndexes(t, repository, "id0", "battleId0", "userId0", 0)
}

func newEvent(eventJSON string) *stream.Event {
	event := &stream.Event{}
	json.Unmarshal([]byte(eventJSON), event)
	return event
}

func setupIndexes(repository *mocks.Repository, commentID string, battleID string, userID string, score float64) {
	repository.SetScore(fmt.Sprintf(battleKeyPattern, battleID), commentID, score)
	repository.SetScore(fmt.Sprintf(authorKeyPattern, userID), commentID, score)
}

func verifyIndexes(t *testing.T, repository *mocks.Repository, commentID string, battleID string, userID string, expectedScore float64) {
	scoreByBattle := repository.GetScore(fmt.Sprintf(battleKeyPattern, battleID), commentID)
	scoreByAuthor := repository.GetScore(fmt.Sprintf(authorKeyPattern, userID), commentID)
	AssertEquals(t, scoreByBattle, expectedScore)
	AssertEquals(t, scoreByAuthor, expectedScore)
}
