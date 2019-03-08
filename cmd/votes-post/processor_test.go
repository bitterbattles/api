package main_test

import (
	"testing"

	. "github.com/bitterbattles/api/cmd/votes-post"
	"github.com/bitterbattles/api/pkg/errors"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
	"github.com/bitterbattles/api/pkg/votes/mocks"
)

func TestProcessorNoBattleId(t *testing.T) {
	testProcessor(t, "", true, http.BadRequest)
}

func TestProcessor(t *testing.T) {
	testProcessor(t, "bbbbbbbbbbbbbattleID", true, http.Created)
}

func testProcessor(t *testing.T, battleID string, isVoteFor bool, expectedStatusCode int) {
	authContext := &api.AuthContext{
		UserID: "userId",
	}
	requestBody := &Request{
		BattleID:  battleID,
		IsVoteFor: isVoteFor,
	}
	input := &api.Input{
		AuthContext: authContext,
		RequestBody: requestBody,
	}
	repository := mocks.NewRepository()
	processor := NewProcessor(repository)
	output, err := processor.Process(input)
	vote := repository.GetLastAdded()
	if expectedStatusCode == http.Created {
		AssertNotNil(t, output)
		AssertNil(t, err)
		AssertEquals(t, output.StatusCode, expectedStatusCode)
		AssertNotNil(t, vote)
		AssertEquals(t, vote.UserID, authContext.UserID)
		AssertEquals(t, vote.BattleID, battleID)
		AssertEquals(t, vote.IsVoteFor, isVoteFor)
	} else {
		AssertNil(t, output)
		AssertNotNil(t, err)
		AssertNil(t, vote)
		httpError, ok := err.(*errors.HTTPError)
		AssertTrue(t, ok)
		AssertEquals(t, httpError.StatusCode(), expectedStatusCode)
	}
}
