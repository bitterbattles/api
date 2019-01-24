package battlesstream_test

import (
	"testing"
	"time"

	"github.com/bitterbattles/api/battles"
	"github.com/bitterbattles/api/battles/mocks"
	. "github.com/bitterbattles/api/common/tests"
	. "github.com/bitterbattles/api/lambdas/battlesstream"
)

func TestHandler(t *testing.T) {
	index := mocks.NewIndex()
	handler := NewHandler(index, nil)
	eventJSON := `{"records":[{"dynamodb":{"NewImage":{"id":{"S":"id0"},"votesFor":{"N":"0"},"votesAgainst":{"N":"0"},"createdOn":{"N":"123"}}}},{"dynamodb":{"NewImage":{"id":{"S":"id1"},"votesFor":{"N":"0"},"votesAgainst":{"N":"0"},"createdOn":{"N":"456"}}}},{"dynamodb":{"OldImage":{"id":{"S":"id0"},"votesFor":{"N":"0"},"votesAgainst":{"N":"0"},"createdOn":{"N":"123"}},"NewImage":{"id":{"S":"id0"},"votesFor":{"N":"10"},"votesAgainst":{"N":"5"},"createdOn":{"N":"123"}}}},{"dynamodb":{"OldImage":{"id":{"S":"id1"},"votesFor":{"N":"0"},"votesAgainst":{"N":"0"},"createdOn":{"N":"456"}},"NewImage":{"id":{"S":"id1"},"votesFor":{"N":"20"},"votesAgainst":{"N":"21"},"createdOn":{"N":"456"}}}}]}`
	responseBytes, err := handler.Invoke(nil, []byte(eventJSON))
	AssertNil(t, responseBytes)
	AssertNil(t, err)
	recencyWeight := float64(time.Now().Unix() / 86400)
	verifyIndexScores(t, index, "id0", 123, recencyWeight+15, recencyWeight+10)
	verifyIndexScores(t, index, "id1", 456, recencyWeight+41, recencyWeight+40)
}

func verifyIndexScores(t *testing.T, index *mocks.Index, id string, expectedRecency float64, expectedPopularity float64, expectedControversy float64) {
	recency := index.GetScore(battles.RecentSort, id)
	popularity := index.GetScore(battles.PopularSort, id)
	controversy := index.GetScore(battles.ControversialSort, id)
	AssertEquals(t, recency, expectedRecency)
	AssertEquals(t, popularity, expectedPopularity)
	AssertEquals(t, controversy, expectedControversy)
}
