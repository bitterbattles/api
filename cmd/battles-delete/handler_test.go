package main_test

import (
	"testing"

	. "github.com/bitterbattles/api/cmd/battles-delete"
	"github.com/bitterbattles/api/pkg/battles"
	battlesMocks "github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/common/http"
	. "github.com/bitterbattles/api/pkg/common/tests"
	ranksMocks "github.com/bitterbattles/api/pkg/ranks/mocks"
)

const sort = battles.RecentSort

func TestHandler(t *testing.T) {
	ranksRepository := ranksMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	id := "id"
	battle := battles.Battle{
		ID: id,
	}
	battlesRepository.Add(battle)
	ranksRepository.SetScore(battles.RecentSort, id, 123)
	ranksRepository.SetScore(battles.PopularSort, id, 456)
	ranksRepository.SetScore(battles.ControversialSort, id, 789)
	pathParams := make(map[string]string)
	pathParams["id"] = id
	request := &http.Request{
		PathParams: pathParams,
	}
	handler := NewHandler(ranksRepository, battlesRepository)
	response, err := handler.Handle(request)
	AssertNil(t, err)
	AssertNotNil(t, response)
	AssertEquals(t, response.StatusCode, http.NoContent)
	foundBattle, _ := battlesRepository.GetByID(id)
	AssertNil(t, foundBattle)
	ranks, _ := ranksRepository.GetRankedBattleIDs(battles.RecentSort, 0, 1)
	AssertEquals(t, len(ranks), 0)
	ranks, _ = ranksRepository.GetRankedBattleIDs(battles.PopularSort, 0, 1)
	AssertEquals(t, len(ranks), 0)
	ranks, _ = ranksRepository.GetRankedBattleIDs(battles.ControversialSort, 0, 1)
	AssertEquals(t, len(ranks), 0)
}
