package main_test

import (
	"encoding/json"
	"testing"

	. "github.com/bitterbattles/api/cmd/votes-post"
	"github.com/bitterbattles/api/pkg/common/errors"
	"github.com/bitterbattles/api/pkg/common/http"
	. "github.com/bitterbattles/api/pkg/common/tests"
	"github.com/bitterbattles/api/pkg/votes/mocks"
)

func TestHandlerNoBattleId(t *testing.T) {
	testHandler(t, "", true, http.BadRequest)
}

func TestHandler(t *testing.T) {
	testHandler(t, "bbbbbbbbbbbbbattleID", true, http.Created)
}

func testHandler(t *testing.T, battleID string, isVoteFor bool, expectedStatusCode int) {
	bodyJSON, _ := json.Marshal(Request{
		BattleID:  battleID,
		IsVoteFor: isVoteFor,
	})
	request := &http.Request{
		Body: string(bodyJSON),
	}
	repository := mocks.NewRepository()
	handler := NewHandler(repository)
	response, err := handler.Handle(request)
	vote := repository.GetLastAdded()
	if expectedStatusCode == http.Created {
		AssertNotNil(t, response)
		AssertNil(t, err)
		AssertEquals(t, response.StatusCode, expectedStatusCode)
		AssertNotNil(t, vote)
		AssertEquals(t, vote.BattleID, battleID)
		AssertEquals(t, vote.IsVoteFor, isVoteFor)
	} else {
		AssertNil(t, response)
		AssertNotNil(t, err)
		AssertNil(t, vote)
		httpError, ok := err.(errors.HTTPError)
		AssertTrue(t, ok)
		AssertEquals(t, httpError.StatusCode, expectedStatusCode)
	}
}
