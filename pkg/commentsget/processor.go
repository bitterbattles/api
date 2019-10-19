package commentsget

import (
	"log"

	"github.com/bitterbattles/api/pkg/comments"
	"github.com/bitterbattles/api/pkg/input"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/users"
)

const (
	pageParam       = "page"
	pageSizeParam   = "pageSize"
	minPage         = 1
	defaultPage     = 1
	minPageSize     = 1
	maxPageSize     = 100
	defaultPageSize = 50
	deletedComment  = "[Deleted]"
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
func CreateResponses(
	commentIDs []string,
	includeDeleted bool,
	commentsRepository comments.RepositoryInterface,
	usersRepository users.RepositoryInterface) ([]*Response, error) {
	responses := make([]*Response, 0, len(commentIDs))
	usernames := make(map[string]string)
	for _, commentID := range commentIDs {
		comment, err := commentsRepository.GetByID(commentID)
		if err != nil {
			return nil, err
		}
		if comment == nil {
			log.Println("Comment with ID", commentID, "is either missing or deleted.")
			continue
		}
		if comment.State == comments.Deleted && !includeDeleted {
			continue
		}
		response, err := createResponseWithUsernameMap(comment, usernames, usersRepository)
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

func createResponseWithUsernameMap(comment *comments.Comment, usernames map[string]string, usersRepository users.RepositoryInterface) (*Response, error) {
	username, err := getUsername(comment.UserID, usernames, usersRepository)
	if err != nil {
		return nil, err
	}
	commentText := comment.Comment
	if comment.State == comments.Deleted {
		commentText = deletedComment
	}
	response := &Response{
		ID:        comment.ID,
		BattleID:  comment.BattleID,
		CreatedOn: comment.CreatedOn,
		Username:  username,
		Comment:   commentText,
	}
	return response, nil
}

func getUsername(userID string, usernames map[string]string, usersRepository users.RepositoryInterface) (string, error) {
	var username string
	if usernames != nil {
		username := usernames[userID]
		if username != "" {
			return username, nil
		}
	}
	user, err := usersRepository.GetByID(userID)
	if err != nil {
		return "", err
	}
	if user != nil {
		username = user.Username
	}
	usernames[userID] = username
	return username, nil
}
