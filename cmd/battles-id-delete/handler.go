package main

import (
	"log"

	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/common/http"
	"github.com/bitterbattles/api/pkg/common/lambda/api"
)

const (
	idParam = "id"
)

// Handler represents a request handler
type Handler struct {
	*api.Handler
	repository battles.RepositoryInterface
}

// NewHandler creates a new Handler instance
func NewHandler(repository battles.RepositoryInterface) *api.Handler {
	handler := Handler{
		repository: repository,
	}
	return api.NewHandler(&handler)
}

// Handle handles a request
func (handler *Handler) Handle(request *http.Request) (*http.Response, error) {
	log.Println(request)
	battleID := request.PathParams[idParam]
	handler.repository.DeleteByID(battleID)
	return http.NewResponseWithStatus(nil, nil, http.NoContent)
}
