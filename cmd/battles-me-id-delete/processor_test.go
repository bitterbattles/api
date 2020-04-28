package main_test

import (
	"testing"

	. "github.com/bitterbattles/api/cmd/battles-me-id-delete"
	"github.com/bitterbattles/api/pkg/battles"
	battlesMocks "github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
)

func TestProcessorAsNotAuthor(t *testing.T) {
	testProcessor(t, true)
}

func TestProcessorAsAuthor(t *testing.T) {
	testProcessor(t, false)
}

func testProcessor(t *testing.T, isBattleAuthor bool) {
	battleID := "battleId"
	userID := "userId"
	battleUserID := "userId"
	if !isBattleAuthor {
		battleUserID = "otherUserId"
	}
	battlesRepository := battlesMocks.NewRepository()
	battle := battles.Battle{
		ID:     battleID,
		UserID: battleUserID,
	}
	battlesRepository.Add(&battle)
	pathParams := make(map[string]string)
	pathParams["id"] = battleID
	authContext := &api.AuthContext{
		UserID: userID,
	}
	input := &api.Input{
		PathParams:  pathParams,
		AuthContext: authContext,
	}
	processor := NewProcessor(battlesRepository)
	output, err := processor.Process(input)
	AssertNil(t, err)
	AssertNotNil(t, output)
	foundBattle, _ := battlesRepository.GetByID(battleID)
	if isBattleAuthor {
		AssertEquals(t, output.StatusCode, http.NoContent)
		AssertNil(t, foundBattle)
	} else {
		AssertEquals(t, output.StatusCode, http.NotFound)
		AssertNotNil(t, foundBattle)
	}
}
