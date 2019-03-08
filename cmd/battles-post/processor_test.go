package main_test

import (
	"testing"

	. "github.com/bitterbattles/api/cmd/battles-post"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/errors"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
)

func TestProcessorTooShortTitle(t *testing.T) {
	testProcessor(t, "", "description", http.BadRequest)
}

func TestProcessorTooLongTitle(t *testing.T) {
	title := "loooooooooooooooooooooooooooooooooooooooooooooooooongtitle"
	testProcessor(t, title, "description", http.BadRequest)
}

func TestProcessorTooShortDescription(t *testing.T) {
	testProcessor(t, "title", "", http.BadRequest)
}

func TestProcessorTooLongDescription(t *testing.T) {
	description := "loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooongdescription"
	testProcessor(t, "title", description, http.BadRequest)
}

func TestProcessorSuccess(t *testing.T) {
	testProcessor(t, "title", "description", http.Created)
}

func testProcessor(t *testing.T, title string, description string, expectedStatusCode int) {
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
	repository := mocks.NewRepository()
	processor := NewProcessor(repository)
	output, err := processor.Process(input)
	battle := repository.GetLastAdded()
	if expectedStatusCode == http.Created {
		AssertNotNil(t, output)
		AssertNil(t, err)
		AssertEquals(t, output.StatusCode, expectedStatusCode)
		AssertNotNil(t, battle)
		AssertEquals(t, battle.UserID, authContext.UserID)
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
