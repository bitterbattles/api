package main

import (
	"log"

	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/input"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/ranks"
	"github.com/bitterbattles/api/pkg/users"
	"github.com/bitterbattles/api/pkg/votes"
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
	defaultUsername = "[Deleted]"
)

// Processor represents a request processor
type Processor struct {
	battlesRepository battles.RepositoryInterface
	ranksRepository   ranks.RepositoryInterface
	usersRepository   users.RepositoryInterface
	votesRepository   votes.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(
	battlesRepository battles.RepositoryInterface,
	ranksRepository ranks.RepositoryInterface,
	usersRepository users.RepositoryInterface,
	votesRepository votes.RepositoryInterface) *Processor {
	return &Processor{
		battlesRepository: battlesRepository,
		ranksRepository:   ranksRepository,
		usersRepository:   usersRepository,
		votesRepository:   votesRepository,
	}
}

// NewRequestBody creates a new request body instance
func (processor *Processor) NewRequestBody() interface{} {
	return nil
}

// Process processes a request
func (processor *Processor) Process(input *api.Input) (*api.Output, error) {
	sort := processor.sanitizeSort(input.QueryParams[sortParam])
	page := processor.sanitizePage(input.QueryParams[pageParam])
	pageSize := processor.sanitizePageSize(input.QueryParams[pageSizeParam])
	ids, err := processor.getIDs(sort, page, pageSize)
	if err != nil {
		return nil, err
	}
	responses := make([]*Response, 0, len(ids))
	for _, id := range ids {
		battle, err := processor.getBattle(id)
		if err != nil {
			return nil, err
		}
		if battle != nil {
			user, err := processor.usersRepository.GetByID(battle.UserID)
			if err != nil {
				return nil, err
			}
			var vote *votes.Vote
			if input.AuthContext != nil && input.AuthContext.UserID != "" {
				vote, err = processor.votesRepository.GetByUserAndBattleIDs(input.AuthContext.UserID, battle.ID)
				if err != nil {
					return nil, err
				}
			}
			responses = append(responses, processor.toResponse(battle, user, vote))
		} else {
			log.Println("Failed to find battle ID", id, "referenced in", sort, "ranking.")
		}
	}
	output := api.NewOutput(responses)
	return output, nil
}

func (processor *Processor) getIDs(sort string, page int, pageSize int) ([]string, error) {
	offset := (page - 1) * pageSize
	limit := pageSize
	ids, err := processor.ranksRepository.GetRankedBattleIDs(sort, offset, limit)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (processor *Processor) getBattle(id string) (*battles.Battle, error) {
	battle, err := processor.battlesRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return battle, nil
}

func (processor *Processor) sanitizeSort(sort string) string {
	rules := input.StringRules{
		ToLower:      true,
		TrimSpace:    true,
		ValidValues:  []string{battles.RecentSort, battles.PopularSort, battles.ControversialSort},
		DefaultValue: battles.RecentSort,
	}
	sort, _ = input.SanitizeString(sort, rules, nil)
	return sort
}

func (processor *Processor) sanitizePage(pageString string) int {
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

func (processor *Processor) sanitizePageSize(pageSizeString string) int {
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

func (processor *Processor) toResponse(battle *battles.Battle, user *users.User, vote *votes.Vote) *Response {
	username := defaultUsername
	if user != nil && user.State == users.Active {
		username = user.DisplayUsername
	}
	var hasVoted bool
	if vote != nil {
		hasVoted = true
	}
	return &Response{
		ID:           battle.ID,
		Username:     username,
		Title:        battle.Title,
		Description:  battle.Description,
		HasVoted:     hasVoted,
		VotesFor:     battle.VotesFor,
		VotesAgainst: battle.VotesAgainst,
		CreatedOn:    battle.CreatedOn,
	}
}
