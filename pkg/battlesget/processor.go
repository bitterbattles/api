package battlesget

import (
	"log"
	"math"

	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/input"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/time"
	"github.com/bitterbattles/api/pkg/users"
	"github.com/bitterbattles/api/pkg/votes"
)

const (
	sortParam         = "sort"
	pageParam         = "page"
	pageSizeParam     = "pageSize"
	recentSort        = "recent"
	popularSort       = "popular"
	controversialSort = "controversial"
	defaultSort       = recentSort
	minPage           = 1
	defaultPage       = 1
	minPageSize       = 1
	maxPageSize       = 100
	defaultPageSize   = 50
	deletedUsername   = "[Deleted]"
)

// GetSort gets and sanitizes the sort param from a GET request
func GetSort(i *api.Input) string {
	sort := i.QueryParams[sortParam]
	rules := input.StringRules{
		ToLower:      true,
		TrimSpace:    true,
		ValidValues:  []string{recentSort, popularSort, controversialSort},
		DefaultValue: defaultSort,
	}
	sort, _ = input.SanitizeString(sort, rules, nil)
	return sort
}

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

// CreateResponses creates a list of GET battles responses
func CreateResponses(
	userID string,
	battleIDs []string,
	battlesRepository battles.RepositoryInterface,
	usersRepository users.RepositoryInterface,
	votesRepository votes.RepositoryInterface) ([]*Response, error) {
	responses := make([]*Response, 0, len(battleIDs))
	usernames := make(map[string]string)
	for _, battleID := range battleIDs {
		response, err := createResponseWithUsernameMap(userID, battleID, usernames, battlesRepository, usersRepository, votesRepository)
		if err != nil {
			return nil, err
		}
		if response == nil {
			log.Println("Battle with ID", battleID, "is either missing or deleted.")
			continue
		}
		responses = append(responses, response)
	}
	return responses, nil
}

// CreateResponse creates a GET battle response
func CreateResponse(
	userID string,
	battleID string,
	battlesRepository battles.RepositoryInterface,
	usersRepository users.RepositoryInterface,
	votesRepository votes.RepositoryInterface) (*Response, error) {
	return createResponseWithUsernameMap(userID, battleID, nil, battlesRepository, usersRepository, votesRepository)
}

func createResponseWithUsernameMap(
	userID string,
	battleID string,
	usernames map[string]string,
	battlesRepository battles.RepositoryInterface,
	usersRepository users.RepositoryInterface,
	votesRepository votes.RepositoryInterface) (*Response, error) {
	battle, err := battlesRepository.GetByID(battleID)
	if err != nil {
		return nil, err
	}
	if battle == nil || battle.State == battles.Deleted {
		return nil, nil
	}
	username, err := getUsername(battle.UserID, usernames, usersRepository)
	if err != nil {
		return nil, err
	}
	canVote, err := getCanVote(userID, battleID, votesRepository)
	if err != nil {
		return nil, err
	}
	response := &Response{
		ID:           battle.ID,
		CreatedOn:    battle.CreatedOn,
		Username:     username,
		Title:        battle.Title,
		Description:  battle.Description,
		CanVote:      canVote,
		VotesFor:     battle.VotesFor,
		VotesAgainst: battle.VotesAgainst,
		Comments:     battle.Comments,
		Verdict:      determineVerdict(battle.CreatedOn, float64(battle.VotesFor), float64(battle.VotesAgainst)),
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
	} else {
		username = deletedUsername
	}
	if usernames != nil {
		usernames[userID] = username
	}
	return username, nil
}

func getCanVote(userID string, battleID string, votesRepository votes.RepositoryInterface) (bool, error) {
	if userID == "" {
		return false, nil
	}
	vote, err := votesRepository.GetByUserAndBattleIDs(userID, battleID)
	if err != nil {
		return false, err
	}
	return (vote == nil), nil
}

func determineVerdict(createdOn int64, votesFor float64, votesAgainst float64) int {
	daysOld := (time.NowUnix() - createdOn) / 86400
	if daysOld < 1 {
		return None
	}
	totalVotes := votesFor + votesAgainst
	deltaVotes := math.Abs(votesFor - votesAgainst)
	if totalVotes == 0 || deltaVotes/totalVotes <= 0.05 {
		return NoDecision
	}
	if votesAgainst > votesFor {
		return Against
	}
	return For
}
