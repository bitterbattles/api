package main_test

import (
	"testing"

	. "github.com/bitterbattles/api/cmd/comments-me-id-delete"
	"github.com/bitterbattles/api/pkg/comments"
	commentsMocks "github.com/bitterbattles/api/pkg/comments/mocks"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
)

func TestProcessorAsNotAuthor(t *testing.T) {
	testProcessor(t, true)
}

func TestProcessorAsAuthor(t *testing.T) {
	testProcessor(t, false)
}

func testProcessor(t *testing.T, isAuthor bool) {
	commentID := "commentId"
	userID := "userId"
	commentUserID := "userId"
	if !isAuthor {
		commentUserID = "otherUserId"
	}
	commentsRepository := commentsMocks.NewRepository()
	comment := comments.Comment{
		ID:     commentID,
		UserID: commentUserID,
	}
	commentsRepository.Add(&comment)
	pathParams := make(map[string]string)
	pathParams["id"] = commentID
	authContext := &api.AuthContext{
		UserID: userID,
	}
	input := &api.Input{
		PathParams:  pathParams,
		AuthContext: authContext,
	}
	processor := NewProcessor(commentsRepository)
	output, err := processor.Process(input)
	AssertNil(t, err)
	AssertNotNil(t, output)
	foundComment, _ := commentsRepository.GetByID(commentID)
	if isAuthor {
		AssertEquals(t, output.StatusCode, http.NoContent)
		AssertNil(t, foundComment)
	} else {
		AssertEquals(t, output.StatusCode, http.NotFound)
		AssertNotNil(t, foundComment)
	}
}
