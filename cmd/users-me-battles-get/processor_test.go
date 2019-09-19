package main_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/bitterbattles/api/cmd/users-me-battles-get"
	"github.com/bitterbattles/api/pkg/battles"
	battlesMocks "github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/http"
	indexMocks "github.com/bitterbattles/api/pkg/index/mocks"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
)

const testSort = "recent"

func TestProcessorBadIndexEntry(t *testing.T) {
	userID := "userId0"
	indexRepository := indexMocks.NewRepository()
	key := fmt.Sprintf("battleIds:forAuthor:%s:%s", userID, testSort)
	indexRepository.SetScore(key, "badId", 0)
	repository := battlesMocks.NewRepository()
	addBattles(indexRepository, repository, userID, testSort, false, 1)
	expectedResponse := `[{"id":"id0","username":"username","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"createdOn":0}]`
	authContext := &api.AuthContext{
		UserID: userID,
	}
	testProcessor(t, indexRepository, repository, authContext, testSort, "1", "2", expectedResponse)
}

func TestProcessorDeletedIndexEntry(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	repository := battlesMocks.NewRepository()
	addBattles(indexRepository, repository, "userId0", testSort, true, 1)
	expectedResponse := `[]`
	authContext := &api.AuthContext{
		UserID: "userId0",
	}
	testProcessor(t, indexRepository, repository, authContext, testSort, "1", "2", expectedResponse)
}

func TestProcessor(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	repository := battlesMocks.NewRepository()
	addBattles(indexRepository, repository, "userId0", testSort, false, 3)
	expectedResponse := `[{"id":"id0","username":"username","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"createdOn":0},{"id":"id1","username":"username","title":"title1","description":"description1","canVote":false,"votesFor":1,"votesAgainst":2,"createdOn":3}]`
	authContext := &api.AuthContext{
		UserID: "userId0",
	}
	testProcessor(t, indexRepository, repository, authContext, testSort, "1", "2", expectedResponse)
}

func addBattles(indexRepository *indexMocks.Repository, repository *battlesMocks.Repository, userID string, sort string, isDeleted bool, count int) {
	key := fmt.Sprintf("battleIds:forAuthor:%s:%s", userID, sort)
	state := battles.Active
	if isDeleted {
		state = battles.Deleted
	}
	for i := 0; i < count; i++ {
		battle := battles.Battle{
			ID:           fmt.Sprintf("id%d", i),
			UserID:       userID,
			Username:     "username",
			Title:        fmt.Sprintf("title%d", i),
			Description:  fmt.Sprintf("description%d", i),
			VotesFor:     i,
			VotesAgainst: i * 2,
			CreatedOn:    int64(i * 3),
			State:        state,
		}
		repository.Add(&battle)
		indexRepository.SetScore(key, battle.ID, float64(i))
	}
}

func testProcessor(t *testing.T, indexRepository *indexMocks.Repository, repository *battlesMocks.Repository, authContext *api.AuthContext, testSort string, page string, pageSize string, expectedResponseBody string) {
	queryParams := make(map[string]string)
	if testSort != "" {
		queryParams["sort"] = testSort
	}
	if page != "" {
		queryParams["page"] = page
	}
	if pageSize != "" {
		queryParams["pageSize"] = pageSize
	}
	input := &api.Input{
		AuthContext: authContext,
		QueryParams: queryParams,
	}
	indexer := battles.NewIndexer(indexRepository)
	processor := NewProcessor(indexer, repository)
	output, err := processor.Process(input)
	AssertNil(t, err)
	AssertNotNil(t, output)
	AssertEquals(t, output.StatusCode, http.OK)
	AssertNotNil(t, output.ResponseBody)
	responseBody, _ := json.Marshal(output.ResponseBody)
	AssertEquals(t, string(responseBody), expectedResponseBody)
}
