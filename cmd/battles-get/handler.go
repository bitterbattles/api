package main

import (
	"log"

	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/common/handlers"
	"github.com/bitterbattles/api/pkg/common/input"
	"github.com/bitterbattles/api/pkg/ranks"
)

const (
	defaultSort     = battles.RecentSort
	minPage         = 1
	defaultPage     = 1
	minPageSize     = 1
	maxPageSize     = 100
	defaultPageSize = 50
)

// Handler represents a request handler
type Handler struct {
	*handlers.APIHandler
	ranksRepository   ranks.RepositoryInterface
	battlesRepository battles.RepositoryInterface
}

// NewHandler creates a new Handler instance
func NewHandler(ranksRepository ranks.RepositoryInterface, battlesRepository battles.RepositoryInterface) *handlers.APIHandler {
	handler := Handler{
		ranksRepository:   ranksRepository,
		battlesRepository: battlesRepository,
	}
	return handlers.NewAPIHandler(handler.Handle)
}

// Handle handles a request
func (handler *Handler) Handle(request *Request) ([]Response, error) {
	sort := handler.sanitizeSort(request.Sort)
	page := handler.sanitizePage(request.Page)
	pageSize := handler.sanitizePageSize(request.PageSize)
	ids, err := handler.getIDs(sort, page, pageSize)
	if err != nil {
		return nil, err
	}
	battles := make([]*battles.Battle, 0, len(ids))
	for _, id := range ids {
		battle, err := handler.getBattle(id)
		if err != nil {
			return nil, err
		}
		if battle != nil {
			battles = append(battles, battle)
		} else {
			log.Println("Failed to find battle ID", id, "referenced in", sort, "ranking.")
		}
	}
	count := len(battles)
	response := make([]Response, count)
	for i := 0; i < count; i++ {
		response[i] = Response(*battles[i])
	}
	return response, nil
}

func (handler *Handler) getIDs(sort string, page int, pageSize int) ([]string, error) {
	offset := (page - 1) * pageSize
	limit := pageSize
	ids, err := handler.ranksRepository.GetRankedBattleIDs(sort, offset, limit)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (handler *Handler) getBattle(id string) (*battles.Battle, error) {
	battle, err := handler.battlesRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return battle, nil
}

func (handler *Handler) sanitizeSort(sort string) string {
	rules := input.StringRules{
		ToLower:      true,
		TrimSpace:    true,
		ValidValues:  []string{battles.RecentSort, battles.PopularSort, battles.ControversialSort},
		DefaultValue: battles.RecentSort,
	}
	sort, _ = input.SanitizeString(sort, rules, nil)
	return sort
}

func (handler *Handler) sanitizePage(page int) int {
	rules := input.IntegerRules{
		EnforceMinValue:    true,
		MinValue:           minPage,
		UseDefaultMinValue: true,
		DefaultMinValue:    defaultPage,
	}
	page, _ = input.SanitizeInteger(page, rules, nil)
	return page
}

func (handler *Handler) sanitizePageSize(pageSize int) int {
	rules := input.IntegerRules{
		EnforceMinValue:    true,
		MinValue:           minPageSize,
		UseDefaultMinValue: true,
		DefaultMinValue:    defaultPageSize,
		EnforceMaxValue:    true,
		MaxValue:           maxPageSize,
		UseDefaultMaxValue: true,
		DefaultMaxValue:    maxPageSize,
	}
	pageSize, _ = input.SanitizeInteger(pageSize, rules, nil)
	return pageSize
}
