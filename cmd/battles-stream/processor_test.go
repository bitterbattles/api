package main_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/bitterbattles/api/cmd/battles-stream"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/index/mocks"
	"github.com/bitterbattles/api/pkg/lambda/stream"
	. "github.com/bitterbattles/api/pkg/tests"
	"github.com/bitterbattles/api/pkg/time"
)

const globalKeyPattern = "battleIds:%s"
const authorKeyPattern = "battleIds:forAuthor:%s:%s"

func TestProcessorNewBattle(t *testing.T) {
	repository := mocks.NewRepository()
	indexer := battles.NewIndexer(repository)
	processor := NewProcessor(indexer)
	eventJSON := `{"records":[{"dynamodb":{"NewImage":{"id":{"S":"id0"},"userId":{"S":"userId0"},"votesFor":{"N":"0"},"votesAgainst":{"N":"0"},"createdOn":{"N":"123"}}}}]}`
	event := newEvent(eventJSON)
	err := processor.Process(event)
	AssertNil(t, err)
	recencyWeight := getRecencyWeight()
	verifyIndexes(t, repository, "id0", "userId0", 123, recencyWeight, recencyWeight)
}

func TestProcessorUpdatedBattle(t *testing.T) {
	repository := mocks.NewRepository()
	indexer := battles.NewIndexer(repository)
	processor := NewProcessor(indexer)
	eventJSON := `{"records":[{"dynamodb":{"OldImage":{"id":{"S":"id0"},"userId":{"S":"userId0"},"votesFor":{"N":"0"},"votesAgainst":{"N":"0"},"createdOn":{"N":"123"}},"NewImage":{"id":{"S":"id0"},"userId":{"S":"userId0"},"votesFor":{"N":"10"},"votesAgainst":{"N":"5"},"createdOn":{"N":"123"}}}},{"dynamodb":{"OldImage":{"id":{"S":"id1"},"userId":{"S":"userId1"},"votesFor":{"N":"0"},"votesAgainst":{"N":"0"},"createdOn":{"N":"456"}},"NewImage":{"id":{"S":"id1"},"userId":{"S":"userId1"},"votesFor":{"N":"20"},"votesAgainst":{"N":"21"},"createdOn":{"N":"456"}}}}]}`
	event := newEvent(eventJSON)
	err := processor.Process(event)
	AssertNil(t, err)
	recencyWeight := getRecencyWeight()
	verifyIndexes(t, repository, "id0", "userId0", 0, recencyWeight+15, recencyWeight+10)
	verifyIndexes(t, repository, "id1", "userId1", 0, recencyWeight+41, recencyWeight+40)
}

func TestProcessorDeletedBattle(t *testing.T) {
	repository := mocks.NewRepository()
	setupIndexes(repository, "id0", "userId0", 1, 2, 3)
	indexer := battles.NewIndexer(repository)
	processor := NewProcessor(indexer)
	eventJSON := `{"records":[{"dynamodb":{"OldImage":{"id":{"S":"id0"},"userId":{"S":"userId0"},"votesFor":{"N":"0"},"votesAgainst":{"N":"0"},"createdOn":{"N":"123"},"state":{"N":"1"}},"NewImage":{"id":{"S":"id0"},"userId":{"S":"userId0"},"votesFor":{"N":"0"},"votesAgainst":{"N":"0"},"createdOn":{"N":"123"},"state":{"N":"2"}}}}]}`
	event := newEvent(eventJSON)
	err := processor.Process(event)
	AssertNil(t, err)
	verifyIndexes(t, repository, "id0", "userId0", 0, 0, 0)
}

func newEvent(eventJSON string) *stream.Event {
	event := &stream.Event{}
	json.Unmarshal([]byte(eventJSON), event)
	return event
}

func setupIndexes(repository *mocks.Repository, battleId string, userID string, recency float64, popularity float64, controversy float64) {
	repository.SetScore(fmt.Sprintf(globalKeyPattern, "recent"), battleId, recency)
	repository.SetScore(fmt.Sprintf(globalKeyPattern, "popular"), battleId, popularity)
	repository.SetScore(fmt.Sprintf(globalKeyPattern, "controversial"), battleId, controversy)
	repository.SetScore(fmt.Sprintf(authorKeyPattern, userID, "recent"), battleId, recency)
	repository.SetScore(fmt.Sprintf(authorKeyPattern, userID, "popular"), battleId, popularity)
	repository.SetScore(fmt.Sprintf(authorKeyPattern, userID, "controversial"), battleId, controversy)
}

func verifyIndexes(t *testing.T, repository *mocks.Repository, battleId string, userID string, expectedRecency float64, expectedPopularity float64, expectedControversy float64) {
	recency := repository.GetScore(fmt.Sprintf(globalKeyPattern, "recent"), battleId)
	popularity := repository.GetScore(fmt.Sprintf(globalKeyPattern, "popular"), battleId)
	controversy := repository.GetScore(fmt.Sprintf(globalKeyPattern, "controversial"), battleId)
	AssertEquals(t, recency, expectedRecency)
	AssertEquals(t, popularity, expectedPopularity)
	AssertEquals(t, controversy, expectedControversy)
	authorRecency := repository.GetScore(fmt.Sprintf(authorKeyPattern, userID, "recent"), battleId)
	authorPopularity := repository.GetScore(fmt.Sprintf(authorKeyPattern, userID, "popular"), battleId)
	authorControversy := repository.GetScore(fmt.Sprintf(authorKeyPattern, userID, "controversial"), battleId)
	AssertEquals(t, authorRecency, expectedRecency)
	AssertEquals(t, authorPopularity, expectedPopularity)
	AssertEquals(t, authorControversy, expectedControversy)
}

func getRecencyWeight() float64 {
	return float64(time.NowUnix() / 86400)
}
