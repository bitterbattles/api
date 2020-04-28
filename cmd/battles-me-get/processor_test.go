package main_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/bitterbattles/api/cmd/battles-me-get"
	"github.com/bitterbattles/api/pkg/battles"
	battlesMocks "github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
	"github.com/bitterbattles/api/pkg/users"
	usersMocks "github.com/bitterbattles/api/pkg/users/mocks"
	"github.com/bitterbattles/api/pkg/votes"
	votesMocks "github.com/bitterbattles/api/pkg/votes/mocks"
)

const testSort = "recent"
const userID = "userId0"

func TestProcessorBadIndexEntry(t *testing.T) {
	battlesIndex := battlesMocks.NewIndex()
	battle := &battles.Battle{
		ID: "badId",
	}
	battlesIndex.Upsert(battle, 0, 0)
	battlesRepository := battlesMocks.NewRepository()
	votesRepository := votesMocks.NewRepository()
	addBattles(battlesIndex, battlesRepository, votesRepository, testSort, false, 1)
	expectedResponse := `[{"id":"id0","createdOn":0,"username":"username0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"comments":0,"verdict":4}]`
	testProcessor(t, battlesIndex, battlesRepository, votesRepository, testSort, "1", "2", expectedResponse)
}

func TestProcessorDeletedIndexEntry(t *testing.T) {
	battlesIndex := battlesMocks.NewIndex()
	battlesRepository := battlesMocks.NewRepository()
	votesRepository := votesMocks.NewRepository()
	addBattles(battlesIndex, battlesRepository, votesRepository, testSort, true, 1)
	expectedResponse := `[]`
	testProcessor(t, battlesIndex, battlesRepository, votesRepository, testSort, "1", "2", expectedResponse)
}

func TestProcessor(t *testing.T) {
	battlesIndex := battlesMocks.NewIndex()
	battlesRepository := battlesMocks.NewRepository()
	votesRepository := votesMocks.NewRepository()
	addBattles(battlesIndex, battlesRepository, votesRepository, testSort, false, 3)
	expectedResponse := `[{"id":"id0","createdOn":0,"username":"username0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"comments":0,"verdict":4},{"id":"id1","createdOn":3,"username":"username0","title":"title1","description":"description1","canVote":false,"votesFor":1,"votesAgainst":2,"comments":3,"verdict":3}]`
	testProcessor(t, battlesIndex, battlesRepository, votesRepository, testSort, "1", "2", expectedResponse)
}

func addBattles(battlesIndex *battlesMocks.Index, battlesRepository *battlesMocks.Repository, votesRepository *votesMocks.Repository, sort string, isDeleted bool, count int) {
	state := battles.Active
	if isDeleted {
		state = battles.Deleted
	}
	for i := 0; i < count; i++ {
		battleID := fmt.Sprintf("id%d", i)
		battle := &battles.Battle{
			ID:           battleID,
			UserID:       userID,
			Title:        fmt.Sprintf("title%d", i),
			Description:  fmt.Sprintf("description%d", i),
			VotesFor:     i,
			VotesAgainst: i * 2,
			Comments:     i * 3,
			CreatedOn:    int64(i * 3),
			State:        state,
		}
		battlesRepository.Add(battle)
		battlesIndex.Upsert(battle, 0, 0)
		vote := &votes.Vote{
			BattleID: battleID,
			UserID:   userID,
		}
		votesRepository.Add(vote)
	}
}

func testProcessor(t *testing.T, battlesIndex *battlesMocks.Index, battlesRepository *battlesMocks.Repository, votesRepository *votesMocks.Repository, testSort string, page string, pageSize string, expectedResponseBody string) {
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
	authContext := &api.AuthContext{
		UserID: "userId0",
	}
	input := &api.Input{
		AuthContext: authContext,
		QueryParams: queryParams,
	}
	usersRepository := usersMocks.NewRepository()
	user := &users.User{
		ID:       "userId0",
		Username: "username0",
	}
	usersRepository.Add(user)
	processor := NewProcessor(battlesIndex, battlesRepository, usersRepository, votesRepository)
	output, err := processor.Process(input)
	AssertNil(t, err)
	AssertNotNil(t, output)
	AssertEquals(t, output.StatusCode, http.OK)
	AssertNotNil(t, output.ResponseBody)
	responseBody, _ := json.Marshal(output.ResponseBody)
	AssertEquals(t, string(responseBody), expectedResponseBody)
}
