package main_test

import (
	"testing"

	. "github.com/bitterbattles/api/cmd/battles-id-comments-post"
	"github.com/bitterbattles/api/pkg/comments"
	commentsMocks "github.com/bitterbattles/api/pkg/comments/mocks"
	"github.com/bitterbattles/api/pkg/errors"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
	"github.com/bitterbattles/api/pkg/users"
	usersMocks "github.com/bitterbattles/api/pkg/users/mocks"
)

func TestProcessorTooShortComment(t *testing.T) {
	testProcessor(t, "", false, http.BadRequest)
}

func TestProcessorTooLongComment(t *testing.T) {
	description := "loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooongcomment"
	testProcessor(t, description, false, http.BadRequest)
}

func TestProcessorMissingUser(t *testing.T) {
	testProcessor(t, "description", true, http.Forbidden)
}

func TestProcessorSuccess(t *testing.T) {
	testProcessor(t, "description", false, http.Created)
}

func testProcessor(t *testing.T, commentText string, missingUser bool, expectedStatusCode int) {
	battleID := "battleId"
	pathParams := make(map[string]string)
	pathParams["id"] = battleID
	user := &users.User{
		ID:       "userId",
		Username: "username",
	}
	usersRepository := usersMocks.NewRepository()
	if !missingUser {
		usersRepository.Add(user)
	}
	authContext := &api.AuthContext{
		UserID: "userId",
	}
	requestBody := &Request{
		Comment: commentText,
	}
	input := &api.Input{
		AuthContext: authContext,
		PathParams:  pathParams,
		RequestBody: requestBody,
	}
	commentsRepository := commentsMocks.NewRepository()
	processor := NewProcessor(commentsRepository, usersRepository)
	output, err := processor.Process(input)
	comment := commentsRepository.GetLastAdded()
	if expectedStatusCode == http.Created {
		AssertNotNil(t, output)
		AssertNil(t, err)
		AssertEquals(t, output.StatusCode, expectedStatusCode)
		AssertNotNil(t, comment)
		AssertEquals(t, comment.BattleID, battleID)
		AssertEquals(t, comment.UserID, user.ID)
		AssertEquals(t, comment.Comment, commentText)
		AssertEquals(t, comment.State, comments.Active)
	} else {
		AssertNil(t, output)
		AssertNotNil(t, err)
		AssertNil(t, comment)
		httpError, ok := err.(*errors.HTTPError)
		AssertTrue(t, ok)
		AssertEquals(t, httpError.StatusCode(), expectedStatusCode)
	}
}
