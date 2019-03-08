package main_test

import (
	"testing"

	. "github.com/bitterbattles/api/cmd/logins-post"
	"github.com/bitterbattles/api/pkg/errors"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/jwt"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
	"github.com/bitterbattles/api/pkg/users"
	"github.com/bitterbattles/api/pkg/users/mocks"
)

func TestProcessorMissingUsername(t *testing.T) {
	testProcessor(t, "", "P@ssw0rd", http.BadRequest)
}

func TestProcessorMissingPassword(t *testing.T) {
	testProcessor(t, "UsErNaMe123", "", http.BadRequest)
}

func TestProcessorUnknownUsername(t *testing.T) {
	testProcessor(t, "unknown", "P@ssw0rd", http.BadRequest)
}

func TestProcessorBadPassword(t *testing.T) {
	testProcessor(t, "UsErNaMe123", "incorrect", http.BadRequest)
}

func TestProcessorSuccess(t *testing.T) {
	testProcessor(t, "UsErNaMe123", "P@ssw0rd", http.OK)
}

func testProcessor(t *testing.T, username string, password string, expectedStatusCode int) {
	requestBody := &Request{
		Username: username,
		Password: password,
	}
	input := &api.Input{
		RequestBody: requestBody,
	}
	repository := mocks.NewRepository()
	user := &users.User{
		ID:           "userId",
		Username:     "UsErNaMe123",
		PasswordHash: "$2a$10$dUIBFokQ2L1iBPs.fIM2r.Xpp4xFSPh5PA9LQkeUPAEwWqMU4Lc7K",
	}
	repository.Add(user)
	tokenSecret := "tokenSecret"
	processor := NewProcessor(repository, tokenSecret)
	output, err := processor.Process(input)
	if expectedStatusCode == http.OK {
		AssertNotNil(t, output)
		AssertNil(t, err)
		AssertEquals(t, output.StatusCode, expectedStatusCode)
		response, ok := output.ResponseBody.(*Response)
		AssertTrue(t, ok)
		token := response.AccessToken
		AssertNotNil(t, token)
		authContext := &api.AuthContext{}
		err = jwt.DecodeHS256(token, tokenSecret, authContext)
		AssertNil(t, err)
		AssertEquals(t, authContext.UserID, user.ID)
	} else {
		AssertNil(t, output)
		AssertNotNil(t, err)
		httpError, ok := err.(*errors.HTTPError)
		AssertTrue(t, ok)
		AssertEquals(t, httpError.StatusCode(), expectedStatusCode)
	}
}
