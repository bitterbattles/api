package main_test

import (
	"encoding/json"
	"testing"

	. "github.com/bitterbattles/api/cmd/battles-post"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/common/errors"
	"github.com/bitterbattles/api/pkg/common/http"
	. "github.com/bitterbattles/api/pkg/common/tests"
)

func TestHandlerTooShortTitle(t *testing.T) {
	testHandler(t, "", "description", http.BadRequest)
}

func TestHandlerTooLongTitle(t *testing.T) {
	title := "loooooooooooooooooooooooooooooooooooooooooooooooooongtitle"
	testHandler(t, title, "description", http.BadRequest)
}

func TestHandlerTooShortDescription(t *testing.T) {
	testHandler(t, "title", "", http.BadRequest)
}

func TestHandlerTooLongDescription(t *testing.T) {
	description := "loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooongdescription"
	testHandler(t, "title", description, http.BadRequest)
}

func TestHandler(t *testing.T) {
	testHandler(t, "title", "description", http.Created)
}

func testHandler(t *testing.T, title string, description string, expectedStatusCode int) {
	body, _ := json.Marshal(Request{
		Title:       title,
		Description: description,
	})
	request := &http.Request{
		Body: string(body),
	}
	repository := mocks.NewRepository()
	handler := NewHandler(repository)
	response, err := handler.Handle(request)
	battle := repository.GetLastAdded()
	if expectedStatusCode == http.Created {
		AssertNotNil(t, response)
		AssertNil(t, err)
		AssertEquals(t, response.StatusCode, expectedStatusCode)
		AssertNotNil(t, battle)
		AssertEquals(t, battle.Title, title)
		AssertEquals(t, battle.Description, description)
		AssertEquals(t, battle.State, battles.Active)
	} else {
		AssertNil(t, response)
		AssertNotNil(t, err)
		AssertNil(t, battle)
		httpError, ok := err.(errors.HTTPError)
		AssertTrue(t, ok)
		AssertEquals(t, httpError.StatusCode, expectedStatusCode)
	}
}
