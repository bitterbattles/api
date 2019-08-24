package main_test

import (
	"testing"

	. "github.com/bitterbattles/api/cmd/users-me-delete"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
	"github.com/bitterbattles/api/pkg/users"
	"github.com/bitterbattles/api/pkg/users/mocks"
)

func TestProcessor(t *testing.T) {
	repository := mocks.NewRepository()
	id := "id"
	user := users.User{
		ID: id,
	}
	repository.Add(&user)
	authContext := &api.AuthContext{
		UserID: id,
	}
	input := &api.Input{
		AuthContext: authContext,
	}
	processor := NewProcessor(repository)
	output, err := processor.Process(input)
	AssertNil(t, err)
	AssertNotNil(t, output)
	AssertEquals(t, output.StatusCode, http.NoContent)
	foundUser, _ := repository.GetByID(id)
	AssertNil(t, foundUser)
}
