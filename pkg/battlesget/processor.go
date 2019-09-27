package battlesget

import (
	"log"
	"math"

	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/input"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/time"
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

// GetSort gets and sanitizes the sort param from a GET request
func GetSort(i *api.Input) string {
	sort := i.QueryParams[sortParam]
	rules := input.StringRules{
		ToLower:      true,
		TrimSpace:    true,
		ValidValues:  []string{battles.RecentSort, battles.PopularSort, battles.ControversialSort},
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
func CreateResponses(userID string, battleIDs []string, repository battles.RepositoryInterface, getCanVote func(string, *battles.Battle) (bool, error)) ([]*Response, error) {
	responses := make([]*Response, 0, len(battleIDs))
	for _, battleID := range battleIDs {
		response, err := CreateResponse(userID, battleID, repository, getCanVote)
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
func CreateResponse(userID string, battleID string, repository battles.RepositoryInterface, getCanVote func(string, *battles.Battle) (bool, error)) (*Response, error) {
	battle, err := repository.GetByID(battleID)
	if err != nil {
		return nil, err
	}
	if battle == nil || battle.State == battles.Deleted {
		return nil, nil
	}
	canVote, err := getCanVote(userID, battle)
	if err != nil {
		return nil, err
	}
	response := &Response{
		ID:           battle.ID,
		CreatedOn:    battle.CreatedOn,
		Username:     battle.Username,
		Title:        battle.Title,
		Description:  battle.Description,
		CanVote:      canVote,
		VotesFor:     battle.VotesFor,
		VotesAgainst: battle.VotesAgainst,
		Verdict:      determineVerdict(battle.CreatedOn, float64(battle.VotesFor), float64(battle.VotesAgainst)),
	}
	return response, nil
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
