package main_test

import (
	"encoding/json"
	"testing"

	. "github.com/bitterbattles/api/cmd/battles-post"
	"github.com/bitterbattles/api/pkg/battles/mocks"
	. "github.com/bitterbattles/api/pkg/common/tests"
)

func TestHandlerTooShortTitle(t *testing.T) {
	expectedResponse := `{"statusCode":400,"message":"Invalid title: Minimum length is 4."}`
	testHandler(t, "", "description", expectedResponse, false)
}

func TestHandlerTooLongTitle(t *testing.T) {
	expectedResponse := `{"statusCode":400,"message":"Invalid title: Maximum length is 50."}`
	title := "loooooooooooooooooooooooooooooooooooooooooooooooooongtitle"
	testHandler(t, title, "description", expectedResponse, false)
}

func TestHandlerTooShortDescription(t *testing.T) {
	expectedResponse := `{"statusCode":400,"message":"Invalid description: Minimum length is 4."}`
	testHandler(t, "title", "", expectedResponse, false)
}

func TestHandlerTooLongDescription(t *testing.T) {
	expectedResponse := `{"statusCode":400,"message":"Invalid description: Maximum length is 500."}`
	description := "loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooongdescription"
	testHandler(t, "title", description, expectedResponse, false)
}

func TestHandler(t *testing.T) {
	testHandler(t, "title", "description", "", true)
}

func testHandler(t *testing.T, title string, description string, expectedResponse string, expectedSuccess bool) {
	repository := mocks.NewRepository()
	handler := NewHandler(repository)
	request := Request{
		Title:       title,
		Description: description,
	}
	requestBytes, _ := json.Marshal(request)
	responseBytes, err := handler.Invoke(nil, requestBytes)
	AssertNil(t, err)
	AssertEquals(t, string(responseBytes), expectedResponse)
	battle := repository.GetLastAdded()
	if expectedSuccess {
		AssertNotNil(t, battle)
		AssertEquals(t, battle.Title, title)
		AssertEquals(t, battle.Description, description)
	} else {
		AssertNil(t, battle)
	}
}
