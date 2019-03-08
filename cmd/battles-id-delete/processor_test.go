package main_test

import (
	"testing"

	. "github.com/bitterbattles/api/cmd/battles-id-delete"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
)

func TestProcessor(t *testing.T) {
	repository := mocks.NewRepository()
	id := "id"
	battle := battles.Battle{
		ID: id,
	}
	repository.Add(&battle)
	pathParams := make(map[string]string)
	pathParams["id"] = id
	input := &api.Input{
		PathParams: pathParams,
	}
	processor := NewProcessor(repository)
	output, err := processor.Process(input)
	AssertNil(t, err)
	AssertNotNil(t, output)
	AssertEquals(t, output.StatusCode, http.NoContent)
	foundBattle, _ := repository.GetByID(id)
	AssertNil(t, foundBattle)
}
