package main_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/bitterbattles/api/cmd/votes-me-battles-get"
	"github.com/bitterbattles/api/pkg/battles"
	battlesMocks "github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/http"
	indexMocks "github.com/bitterbattles/api/pkg/index/mocks"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
	"github.com/bitterbattles/api/pkg/users"
	usersMocks "github.com/bitterbattles/api/pkg/users/mocks"
)

func TestProcessor(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, "userId0", 3)
	expectedResponse := `[{"id":"id0","username":"UserID0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"createdOn":0},{"id":"id1","username":"UserID0","title":"title1","description":"description1","canVote":false,"votesFor":1,"votesAgainst":2,"createdOn":3}]`
	authContext := &api.AuthContext{
		UserID: "userId0",
	}
	testProcessor(t, indexRepository, battlesRepository, authContext, "1", "2", expectedResponse)
}

func addBattles(indexRepository *indexMocks.Repository, battlesRepository *battlesMocks.Repository, userID string, count int) {
	key := fmt.Sprintf("battleIds:forVoter:%s", userID)
	for i := 0; i < count; i++ {
		battle := battles.Battle{
			ID:           fmt.Sprintf("id%d", i),
			UserID:       userID,
			Title:        fmt.Sprintf("title%d", i),
			Description:  fmt.Sprintf("description%d", i),
			VotesFor:     i,
			VotesAgainst: i * 2,
			CreatedOn:    int64(i * 3),
			State:        battles.Active,
		}
		battlesRepository.Add(&battle)
		indexRepository.SetScore(key, battle.ID, float64(i))
	}
}

func testProcessor(t *testing.T, indexRepository *indexMocks.Repository, battlesRepository *battlesMocks.Repository, authContext *api.AuthContext, page string, pageSize string, expectedResponseBody string) {
	usersRepository := usersMocks.NewRepository()
	usersRepository.Add(&users.User{
		ID:              "userId0",
		DisplayUsername: "UserID0",
		State:           users.Active,
	})
	queryParams := make(map[string]string)
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
	processor := NewProcessor(indexer, battlesRepository, usersRepository)
	output, err := processor.Process(input)
	AssertNil(t, err)
	AssertNotNil(t, output)
	AssertEquals(t, output.StatusCode, http.OK)
	AssertNotNil(t, output.ResponseBody)
	responseBody, _ := json.Marshal(output.ResponseBody)
	AssertEquals(t, string(responseBody), expectedResponseBody)
}
