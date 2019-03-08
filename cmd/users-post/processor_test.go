package main_test

import (
	"strings"
	"testing"

	. "github.com/bitterbattles/api/cmd/users-post"
	"github.com/bitterbattles/api/pkg/crypto"
	"github.com/bitterbattles/api/pkg/errors"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
	"github.com/bitterbattles/api/pkg/users"
	"github.com/bitterbattles/api/pkg/users/mocks"
)

func TestProcessorTooShortUsername(t *testing.T) {
	testProcessor(t, "a", "P@ssw0rd", http.BadRequest)
}

func TestProcessorTooLongUsername(t *testing.T) {
	testProcessor(t, "toooooooooooooooolong", "P@ssw0rd", http.BadRequest)
}

func TestProcessorInvalidUsername(t *testing.T) {
	testProcessor(t, "us3rn@me", "P@ssw0rd", http.BadRequest)
}

func TestProcessorTooShortPassword(t *testing.T) {
	testProcessor(t, "UsErNaMe123", "P@ss", http.BadRequest)
}

func TestProcessorTooLongPassword(t *testing.T) {
	testProcessor(t, "UsErNaMe123", "P@ssw00000000000000000000rd", http.BadRequest)
}

func TestProcessorInvalidPassword(t *testing.T) {
	testProcessor(t, "UsErNaMe123", "password", http.BadRequest)
}

func TestProcessorUsernameTaken(t *testing.T) {
	testProcessor(t, "ExistingUsername", "P@ssw0rd", http.BadRequest)
}

func TestProcessorSuccess(t *testing.T) {
	testProcessor(t, "UsErNaMe123", "P@ssw0rd", http.Created)
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
	existingUser := &users.User{
		Username: "existingusername",
	}
	repository.Add(existingUser)
	processor := NewProcessor(repository)
	output, err := processor.Process(input)
	user := repository.GetLastAdded()
	if expectedStatusCode == http.Created {
		AssertNotNil(t, output)
		AssertNil(t, err)
		AssertEquals(t, output.StatusCode, expectedStatusCode)
		AssertNotNil(t, user)
		AssertEquals(t, user.Username, strings.ToLower(username))
		AssertEquals(t, user.DisplayUsername, username)
		match := crypto.VerifyPasswordHash(password, user.PasswordHash)
		AssertTrue(t, match)
		AssertEquals(t, user.State, users.Active)
	} else {
		AssertNil(t, output)
		AssertNotNil(t, err)
		AssertEquals(t, user, existingUser)
		httpError, ok := err.(*errors.HTTPError)
		AssertTrue(t, ok)
		AssertEquals(t, httpError.StatusCode(), expectedStatusCode)
	}
}
