package main_test

import (
	"fmt"
	"testing"

	. "github.com/bitterbattles/api/cmd/battles-me-id-delete"
	"github.com/bitterbattles/api/pkg/battles"
	battlesMocks "github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/http"
	indexMocks "github.com/bitterbattles/api/pkg/index/mocks"
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
	indexRepository := indexMocks.NewRepository()
	if isBattleAuthor {
		key := fmt.Sprintf("battleIds:forAuthor:%s:recent", userID)
		indexRepository.SetScore(key, battleID, 1)
	}
	indexer := battles.NewIndexer(indexRepository)
	battlesRepository := battlesMocks.NewRepository()
	battle := battles.Battle{
		ID: battleID,
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
	processor := NewProcessor(indexer, battlesRepository)
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
