package main_test

import (
	"testing"
	"time"

	. "github.com/bitterbattles/api/cmd/battles-stream"
	"github.com/bitterbattles/api/pkg/battles"
	. "github.com/bitterbattles/api/pkg/common/tests"
	"github.com/bitterbattles/api/pkg/ranks/mocks"
)

func TestHandlerNewBattle(t *testing.T) {
	repository := mocks.NewRepository()
	handler := NewHandler(repository)
	eventJSON := `{"records":[{"dynamodb":{"NewImage":{"id":{"S":"id0"},"votesFor":{"N":"0"},"votesAgainst":{"N":"0"},"createdOn":{"N":"123"}}}}]}`
	responseBytes, err := handler.Invoke(nil, []byte(eventJSON))
	AssertNil(t, responseBytes)
	AssertNil(t, err)
	recencyWeight := float64(time.Now().Unix() / 86400)
	verifyRankScores(t, repository, "id0", 123, recencyWeight, recencyWeight)
}

func TestHandlerUpdatedBattle(t *testing.T) {
	repository := mocks.NewRepository()
	handler := NewHandler(repository)
	eventJSON := `{"records":[{"dynamodb":{"NewImage":{"id":{"S":"id0"},"votesFor":{"N":"0"},"votesAgainst":{"N":"0"},"createdOn":{"N":"123"}}}},{"dynamodb":{"NewImage":{"id":{"S":"id1"},"votesFor":{"N":"0"},"votesAgainst":{"N":"0"},"createdOn":{"N":"456"}}}},{"dynamodb":{"OldImage":{"id":{"S":"id0"},"votesFor":{"N":"0"},"votesAgainst":{"N":"0"},"createdOn":{"N":"123"}},"NewImage":{"id":{"S":"id0"},"votesFor":{"N":"10"},"votesAgainst":{"N":"5"},"createdOn":{"N":"123"}}}},{"dynamodb":{"OldImage":{"id":{"S":"id1"},"votesFor":{"N":"0"},"votesAgainst":{"N":"0"},"createdOn":{"N":"456"}},"NewImage":{"id":{"S":"id1"},"votesFor":{"N":"20"},"votesAgainst":{"N":"21"},"createdOn":{"N":"456"}}}}]}`
	responseBytes, err := handler.Invoke(nil, []byte(eventJSON))
	AssertNil(t, responseBytes)
	AssertNil(t, err)
	recencyWeight := float64(time.Now().Unix() / 86400)
	verifyRankScores(t, repository, "id0", 123, recencyWeight+15, recencyWeight+10)
	verifyRankScores(t, repository, "id1", 456, recencyWeight+41, recencyWeight+40)
}

func verifyRankScores(t *testing.T, repository *mocks.Repository, id string, expectedRecency float64, expectedPopularity float64, expectedControversy float64) {
	recency := repository.GetScore(battles.RecentSort, id)
	popularity := repository.GetScore(battles.PopularSort, id)
	controversy := repository.GetScore(battles.ControversialSort, id)
	AssertEquals(t, recency, expectedRecency)
	AssertEquals(t, popularity, expectedPopularity)
	AssertEquals(t, controversy, expectedControversy)
}
