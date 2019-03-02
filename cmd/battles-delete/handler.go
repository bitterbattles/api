package main

import (
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/common/http"
	"github.com/bitterbattles/api/pkg/common/lambda/api"
	"github.com/bitterbattles/api/pkg/ranks"
)

const (
	idParam = "id"
)

// Handler represents a request handler
type Handler struct {
	*api.Handler
	ranksRepository   ranks.RepositoryInterface
	battlesRepository battles.RepositoryInterface
}

// NewHandler creates a new Handler instance
func NewHandler(ranksRepository ranks.RepositoryInterface, battlesRepository battles.RepositoryInterface) *api.Handler {
	handler := Handler{
		ranksRepository:   ranksRepository,
		battlesRepository: battlesRepository,
	}
	return api.NewHandler(&handler)
}

// Handle handles a request
func (handler *Handler) Handle(request *http.Request) (*http.Response, error) {
	battleID := request.PathParams[idParam]
	handler.battlesRepository.DeleteByID(battleID)
	handler.ranksRepository.DeleteByBattleID(battles.RecentSort, battleID)
	handler.ranksRepository.DeleteByBattleID(battles.PopularSort, battleID)
	handler.ranksRepository.DeleteByBattleID(battles.ControversialSort, battleID)
	return http.NewResponseWithStatus(nil, nil, http.NoContent)
}
