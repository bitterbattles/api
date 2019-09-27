package main

import (
	"log"

	"github.com/bitterbattles/api/pkg/comments"
	"github.com/bitterbattles/api/pkg/input"
	"github.com/bitterbattles/api/pkg/lambda/api"
)

const (
	idParam         = "id"
	pageParam       = "page"
	pageSizeParam   = "pageSize"
	minPage         = 1
	defaultPage     = 1
	minPageSize     = 1
	maxPageSize     = 100
	defaultPageSize = 50
)

// Processor represents a request processor
type Processor struct {
	indexer    *comments.Indexer
	repository comments.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(indexer *comments.Indexer, repository comments.RepositoryInterface) *Processor {
	return &Processor{
		indexer:    indexer,
		repository: repository,
	}
}

// NewRequestBody creates a new request body instance
func (processor *Processor) NewRequestBody() interface{} {
	return nil
}

// Process processes a request
func (processor *Processor) Process(input *api.Input) (*api.Output, error) {
	battleID := input.PathParams[idParam]
	page := getPage(input)
	pageSize := getPageSize(input)
	commentIDs, err := processor.indexer.GetByBattleID(battleID, page, pageSize)
	if err != nil {
		return nil, err
	}
	responses := make([]*Response, 0, len(commentIDs))
	for _, commentID := range commentIDs {
		comment, err := processor.repository.GetByID(commentID)
		if err != nil {
			return nil, err
		}
		if comment == nil {
			log.Println("Failed to find comment with ID", commentID, "when fetching comments for battle with ID ", battleID, ".")
			continue
		}
		if comment.State == comments.Deleted {
			continue
		}
		response := &Response{
			ID:        comment.ID,
			CreatedOn: comment.CreatedOn,
			BattleID:  comment.BattleID,
			Username:  comment.Username,
			Comment:   comment.Comment,
		}
		responses = append(responses, response)
	}
	output := api.NewOutput(responses)
	return output, nil
}

func getPage(i *api.Input) int {
	pageString := i.QueryParams[pageParam]
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

func getPageSize(i *api.Input) int {
	pageSizeString := i.QueryParams[pageSizeParam]
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
