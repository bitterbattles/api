package main_test

import (
	"testing"

	. "github.com/bitterbattles/api/cmd/battles-post"
	"github.com/bitterbattles/api/pkg/battles"
	battlesMocks "github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/errors"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
	"github.com/bitterbattles/api/pkg/users"
	usersMocks "github.com/bitterbattles/api/pkg/users/mocks"
)

func TestProcessorTooShortTitle(t *testing.T) {
	testProcessor(t, "", "description", false, http.BadRequest)
}

func TestProcessorTooLongTitle(t *testing.T) {
	title := "loooooooooooooooooooooooooooooooooooooooooooooooooongtitle"
	testProcessor(t, title, "description", false, http.BadRequest)
}

func TestProcessorTooShortDescription(t *testing.T) {
	testProcessor(t, "title", "", false, http.BadRequest)
}

func TestProcessorTooLongDescription(t *testing.T) {
	description := "loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooongdescription"
	testProcessor(t, "title", description, false, http.BadRequest)
}

func TestProcessorMissingUser(t *testing.T) {
	testProcessor(t, "title", "description", true, http.Forbidden)
}

func TestProcessorSuccess(t *testing.T) {
	testProcessor(t, "title", "description", false, http.Created)
}

func testProcessor(t *testing.T, title string, description string, missingUser bool, expectedStatusCode int) {
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
		Title:       title,
		Description: description,
	}
	input := &api.Input{
		AuthContext: authContext,
		RequestBody: requestBody,
	}
	battlesRepository := battlesMocks.NewRepository()
	processor := NewProcessor(battlesRepository, usersRepository)
	output, err := processor.Process(input)
	battle := battlesRepository.GetLastAdded()
	if expectedStatusCode == http.Created {
		AssertNotNil(t, output)
		AssertNil(t, err)
		AssertEquals(t, output.StatusCode, expectedStatusCode)
		AssertNotNil(t, battle)
		AssertEquals(t, battle.UserID, user.ID)
		AssertEquals(t, battle.Title, title)
		AssertEquals(t, battle.Description, description)
		AssertEquals(t, battle.State, battles.Active)
	} else {
		AssertNil(t, output)
		AssertNotNil(t, err)
		AssertNil(t, battle)
		httpError, ok := err.(*errors.HTTPError)
		AssertTrue(t, ok)
		AssertEquals(t, httpError.StatusCode(), expectedStatusCode)
	}
}
