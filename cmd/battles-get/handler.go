package main

import (
	"log"

	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/common/http"
	"github.com/bitterbattles/api/pkg/common/input"
	"github.com/bitterbattles/api/pkg/common/lambda/api"
	"github.com/bitterbattles/api/pkg/ranks"
)

const (
	sortParam       = "sort"
	pageParam       = "page"
	pageSizeParam   = "pageSize"
	defaultSort     = battles.RecentSort
	minPage         = 1
	defaultPage     = 1
	minPageSize     = 1
	maxPageSize     = 100
	defaultPageSize = 50
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
	sort := handler.sanitizeSort(request.QueryParams[sortParam])
	page := handler.sanitizePage(request.QueryParams[pageParam])
	pageSize := handler.sanitizePageSize(request.QueryParams[pageSizeParam])
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
	responseBody := make([]Response, count)
	for i := 0; i < count; i++ {
		responseBody[i] = Response(*battles[i])
	}
	// TODO: Pagination headers
	return http.NewResponse(responseBody, nil)
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

func (handler *Handler) sanitizePage(pageString string) int {
	rules := input.IntegerRules{
		EnforceNumericString: false,
		EnforceMinValue:      true,
		MinValue:             minPage,
		UseDefaultMinValue:   true,
		DefaultMinValue:      defaultPage,
	}
	page, _ := input.SanitizeIntegerFromString(pageString, rules, nil)
	return page
}

func (handler *Handler) sanitizePageSize(pageSizeString string) int {
	rules := input.IntegerRules{
		EnforceNumericString: false,
		EnforceMinValue:      true,
		MinValue:             minPageSize,
		UseDefaultMinValue:   true,
		DefaultMinValue:      defaultPageSize,
		EnforceMaxValue:      true,
		MaxValue:             maxPageSize,
		UseDefaultMaxValue:   true,
		DefaultMaxValue:      maxPageSize,
	}
	pageSize, _ := input.SanitizeIntegerFromString(pageSizeString, rules, nil)
	return pageSize
}
