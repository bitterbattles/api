package getbattles

import (
	"github.com/bitterbattles/api/battles"
	"github.com/bitterbattles/api/common/handlers"
	"github.com/bitterbattles/api/common/input"
	"github.com/bitterbattles/api/common/loggers"
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
	index      battles.IndexInterface
	repository battles.RepositoryInterface
}

// NewHandler creates a new Handler instance
func NewHandler(index battles.IndexInterface, repository battles.RepositoryInterface, logger loggers.LoggerInterface) *handlers.APIHandler {
	handler := Handler{
		index:      index,
		repository: repository,
	}
	return handlers.NewAPIHandler(handler.Handle, logger)
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
			// TODO: Log
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
	ids, err := handler.index.GetRange(sort, offset, limit)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (handler *Handler) getBattle(id string) (*battles.Battle, error) {
	battle, err := handler.repository.GetByID(id)
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
