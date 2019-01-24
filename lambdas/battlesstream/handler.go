package battlesstream

import (
	"math"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bitterbattles/api/battles"
	"github.com/bitterbattles/api/common/handlers"
)

// Handler represents a stream handler
type Handler struct {
	index battles.IndexInterface
}

// NewHandler creates a new Handler instance
func NewHandler(index battles.IndexInterface) *handlers.StreamHandler {
	handler := Handler{
		index: index,
	}
	return handlers.NewStreamHandler(&handler)
}

// TODO: Log errors

// Handle handles a DynamoDB event
func (handler *Handler) Handle(event *handlers.DynamoEvent) error {
	changes := make(map[string]*change)
	for _, record := range event.Records {
		handler.captureChange(&record, changes)
	}
	for battleID, change := range changes {
		handler.processChange(battleID, change)
	}
	return nil
}

func (handler *Handler) captureChange(record *handlers.DynamoEventRecord, changes map[string]*change) {
	var err error
	oldBattle := battles.Battle{}
	err = dynamodbattribute.UnmarshalMap(record.Change.OldImage, &oldBattle)
	if err != nil {
		return
	}
	newBattle := battles.Battle{}
	err = dynamodbattribute.UnmarshalMap(record.Change.NewImage, &newBattle)
	if err != nil {
		return
	}
	c := changes[newBattle.ID]
	if c == nil {
		c = &change{}
	}
	if newBattle.CreatedOn != oldBattle.CreatedOn {
		c.createdOnChanged = true
		c.newCreatedOn = newBattle.CreatedOn
	}
	if newBattle.VotesFor != oldBattle.VotesFor || newBattle.VotesAgainst != oldBattle.VotesAgainst {
		c.votesChanged = true
		c.newVotesFor = newBattle.VotesFor
		c.newVotesAgainst = newBattle.VotesAgainst
	}
	changes[newBattle.ID] = c
}

func (handler *Handler) processChange(battleID string, change *change) {
	if change.createdOnChanged {
		handler.index.Set(battles.RecentSort, battleID, handler.calculateRecency(change.newCreatedOn))
	}
	if change.votesChanged {
		handler.index.Set(battles.PopularSort, battleID, handler.calculatePopularity(change.newVotesFor, change.newVotesAgainst))
		handler.index.Set(battles.ControversialSort, battleID, handler.calculateControversy(change.newVotesFor, change.newVotesAgainst))
	}
}

func (handler *Handler) calculateRecency(createdOn int64) float64 {
	return float64(createdOn) // TODO
}

func (handler *Handler) calculatePopularity(votesFor int, votesAgainst int) float64 {
	return float64(votesFor + votesAgainst) // TODO
}

func (handler *Handler) calculateControversy(votesFor int, votesAgainst int) float64 {
	return float64((votesFor + votesAgainst)) - math.Abs(float64(votesFor-votesAgainst)) // TODO
}
