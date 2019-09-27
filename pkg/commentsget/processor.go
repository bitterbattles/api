package commentsget

import (
	"log"

	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/comments"
	"github.com/bitterbattles/api/pkg/input"
	"github.com/bitterbattles/api/pkg/lambda/api"
)

const (
	pageParam       = "page"
	pageSizeParam   = "pageSize"
	minPage         = 1
	defaultPage     = 1
	minPageSize     = 1
	maxPageSize     = 100
	defaultPageSize = 50
)

// GetPage gets and sanitizes the page param from a GET request
func GetPage(i *api.Input) int {
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

// GetPageSize gets and sanitizes the page size param from a GET request
func GetPageSize(i *api.Input) int {
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

// CreateResponses creates a list of GET comments responses
func CreateResponses(commentIDs []string, repository comments.RepositoryInterface) ([]*Response, error) {
	responses := make([]*Response, 0, len(commentIDs))
	for _, commentID := range commentIDs {
		response, err := CreateResponse(commentID, repository)
		if err != nil {
			return nil, err
		}
		if response == nil {
			log.Println("Comment with ID", commentID, "is either missing or deleted.")
			continue
		}
		responses = append(responses, response)
	}
	return responses, nil
}

// CreateResponse creates a GET comment response
func CreateResponse(commentID string, repository comments.RepositoryInterface) (*Response, error) {
	comment, err := repository.GetByID(commentID)
	if err != nil {
		return nil, err
	}
	if comment == nil || comment.State == battles.Deleted {
		return nil, nil
	}
	response := &Response{
		ID:        comment.ID,
		CreatedOn: comment.CreatedOn,
		BattleID:  comment.BattleID,
		Username:  comment.Username,
		Comment:   comment.Comment,
	}
	return response, nil
}
