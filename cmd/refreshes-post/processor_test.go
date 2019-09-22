package main_test

import (
	"testing"

	. "github.com/bitterbattles/api/cmd/refreshes-post"
	"github.com/bitterbattles/api/pkg/errors"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/jwt"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/loginspost"
	. "github.com/bitterbattles/api/pkg/tests"
	"github.com/bitterbattles/api/pkg/users"
	"github.com/bitterbattles/api/pkg/users/mocks"
)

func TestProcessorBadToken(t *testing.T) {
	testProcessor(t, "invalidToken", http.BadRequest)
}

func TestProcessorMissingUser(t *testing.T) {
	testProcessor(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJtaXNzaW5nVXNlcklkIiwiaWF0IjoxNTY4ODIyODUxLCJleHAiOjI4MzExMjY4NTF9.mQNO1MgG4Cjl1bJDsEeUt7HxICLJTMlP7jIcDfFCpos", http.BadRequest)
}

func TestProcessorExpiredToken(t *testing.T) {
	testProcessor(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VySWQiLCJpYXQiOjE1Njg4MjI4NTEsImV4cCI6MTU2ODgyMjg1MX0.DiA8v2pyB6G40F0d8mJ7nGpGkExy2tdLJd8imreodlo", http.BadRequest)
}

func TestProcessorSuccess(t *testing.T) {
	testProcessor(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VySWQiLCJpYXQiOjE1Njg4MjI4NTEsImV4cCI6MjgzMTEyNjg1MX0.pCqOafVCwQ-9X3u7Lr30UL2fDrOxuPV5AiiA6L2uE5c", http.OK)
}

func testProcessor(t *testing.T, refreshToken string, expectedStatusCode int) {
	requestBody := &Request{
		RefreshToken: refreshToken,
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
	accessTokenSecret := "accessTokenSecret"
	refreshTokenSecret := "refreshTokenSecret"
	processor := NewProcessor(repository, accessTokenSecret, refreshTokenSecret)
	output, err := processor.Process(input)
	if expectedStatusCode == http.OK {
		AssertNotNil(t, output)
		AssertNil(t, err)
		AssertEquals(t, output.StatusCode, expectedStatusCode)
		response, ok := output.ResponseBody.(*loginspost.Response)
		AssertTrue(t, ok)
		accessToken := response.AccessToken
		AssertNotNil(t, accessToken)
		accessAuthContext := &api.AuthContext{}
		err = jwt.DecodeHS256(accessToken, accessTokenSecret, accessAuthContext)
		AssertNil(t, err)
		AssertEquals(t, accessAuthContext.UserID, user.ID)
		AssertEquals(t, response.AccessExpiresIn, 3600)
		refreshToken := response.RefreshToken
		AssertNotNil(t, refreshToken)
		refreshAuthContext := &api.AuthContext{}
		err = jwt.DecodeHS256(refreshToken, refreshTokenSecret, refreshAuthContext)
		AssertNil(t, err)
		AssertEquals(t, refreshAuthContext.UserID, user.ID)
		AssertEquals(t, response.RefreshExpiresIn, 15768000)
	} else {
		AssertNil(t, output)
		AssertNotNil(t, err)
		httpError, ok := err.(*errors.HTTPError)
		AssertTrue(t, ok)
		AssertEquals(t, httpError.StatusCode(), expectedStatusCode)
	}
}
