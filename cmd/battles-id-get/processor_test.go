package main_test

import (
	"encoding/json"
	"testing"

	. "github.com/bitterbattles/api/cmd/battles-id-get"
	"github.com/bitterbattles/api/pkg/battles"
	battlesMocks "github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/errors"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
	"github.com/bitterbattles/api/pkg/users"
	usersMocks "github.com/bitterbattles/api/pkg/users/mocks"
	"github.com/bitterbattles/api/pkg/votes"
	votesMocks "github.com/bitterbattles/api/pkg/votes/mocks"
)

func TestProcessorNotFoundDeleted(t *testing.T) {
	testProcessor(t, "id1", "userId0", http.NotFound, "")
}

func TestProcessorNotFoundMissing(t *testing.T) {
	testProcessor(t, "missing", "userId0", http.NotFound, "")
}

func TestProcessorCanVoteTrue(t *testing.T) {
	expectedResponse := `{"id":"id0","createdOn":2,"username":"username0","title":"title0","description":"description0","canVote":true,"votesFor":0,"votesAgainst":1,"comments":3,"verdict":3}`
	testProcessor(t, "id0", "userId1", http.OK, expectedResponse)
}

func TestProcessorCanVoteFalse(t *testing.T) {
	expectedResponse := `{"id":"id0","createdOn":2,"username":"username0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":1,"comments":3,"verdict":3}`
	testProcessor(t, "id0", "userId0", http.OK, expectedResponse)
}

func testProcessor(t *testing.T, battleID string, userID string, expectedStatusCode int, expectedResponseBody string) {
	battlesRepository := battlesMocks.NewRepository()
	battle := battles.Battle{
		ID:           "id0",
		UserID:       "userId0",
		Title:        "title0",
		Description:  "description0",
		VotesFor:     0,
		VotesAgainst: 1,
		Comments:     3,
		CreatedOn:    2,
		State:        battles.Active,
	}
	battlesRepository.Add(&battle)
	usersRepository := usersMocks.NewRepository()
	user := &users.User{
		ID:       "userId0",
		Username: "username0",
	}
	usersRepository.Add(user)
	votesRepository := votesMocks.NewRepository()
	votesRepository.Add(&votes.Vote{
		UserID:   "userId0",
		BattleID: "id0",
	})
	pathParams := make(map[string]string)
	pathParams["id"] = battleID
	authContext := &api.AuthContext{
		UserID: userID,
	}
	input := &api.Input{
		AuthContext: authContext,
		PathParams:  pathParams,
	}
	processor := NewProcessor(battlesRepository, usersRepository, votesRepository)
	output, err := processor.Process(input)
	if expectedStatusCode == http.OK {
		AssertNil(t, err)
		AssertNotNil(t, output)
		AssertEquals(t, output.StatusCode, http.OK)
		AssertNotNil(t, output.ResponseBody)
		responseBody, _ := json.Marshal(output.ResponseBody)
		AssertEquals(t, string(responseBody), expectedResponseBody)
	} else {
		AssertNil(t, output)
		AssertNotNil(t, err)
		httpError, ok := err.(*errors.HTTPError)
		AssertTrue(t, ok)
		AssertEquals(t, httpError.StatusCode(), expectedStatusCode)
	}
}
