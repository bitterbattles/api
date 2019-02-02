package main_test

import (
	"encoding/json"
	"testing"

	. "github.com/bitterbattles/api/cmd/votes-post"
	. "github.com/bitterbattles/api/pkg/common/tests"
	"github.com/bitterbattles/api/pkg/votes/mocks"
)

func TestHandlerNoBattleId(t *testing.T) {
	expectedResponse := `{"statusCode":400,"message":"Invalid Battle ID: Length must be exactly 20."}`
	testHandler(t, "", true, expectedResponse, false)
}

func TestHandler(t *testing.T) {
	testHandler(t, "bbbbbbbbbbbbbattleID", true, "", true)
}

func testHandler(t *testing.T, battleID string, isVoteFor bool, expectedResponse string, expectedSuccess bool) {
	repository := mocks.NewRepository()
	handler := NewHandler(repository)
	request := Request{
		BattleID:  battleID,
		IsVoteFor: isVoteFor,
	}
	requestBytes, _ := json.Marshal(request)
	responseBytes, err := handler.Invoke(nil, requestBytes)
	AssertNil(t, err)
	AssertEquals(t, string(responseBytes), expectedResponse)
	vote := repository.GetLastAdded()
	if expectedSuccess {
		AssertNotNil(t, vote)
		AssertEquals(t, vote.BattleID, battleID)
		AssertEquals(t, vote.IsVoteFor, isVoteFor)
	} else {
		AssertNil(t, vote)
	}
}
