package main_test

import (
	"testing"

	. "github.com/bitterbattles/api/cmd/battles-id-delete"
	"github.com/bitterbattles/api/pkg/battles"
	battlesMocks "github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/common/http"
	. "github.com/bitterbattles/api/pkg/common/tests"
)

func TestHandler(t *testing.T) {
	repository := battlesMocks.NewRepository()
	id := "id"
	battle := battles.Battle{
		ID: id,
	}
	repository.Add(battle)
	pathParams := make(map[string]string)
	pathParams["id"] = id
	request := &http.Request{
		PathParams: pathParams,
	}
	handler := NewHandler(repository)
	response, err := handler.Handle(request)
	AssertNil(t, err)
	AssertNotNil(t, response)
	AssertEquals(t, response.StatusCode, http.NoContent)
	foundBattle, _ := repository.GetByID(id)
	AssertNil(t, foundBattle)
}
