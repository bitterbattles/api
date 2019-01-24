package battlesstream

import (
	"math"
	"time"

	"github.com/bitterbattles/api/common/loggers"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bitterbattles/api/battles"
	"github.com/bitterbattles/api/common/handlers"
)

// Handler represents a stream handler
type Handler struct {
	index  battles.IndexInterface
	logger loggers.LoggerInterface
}

// NewHandler creates a new Handler instance
func NewHandler(index battles.IndexInterface, logger loggers.LoggerInterface) *handlers.StreamHandler {
	handler := Handler{
		index:  index,
		logger: logger,
	}
	return handlers.NewStreamHandler(&handler)
}

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
		handler.logger.Error("Failed to unmarshal old image in DynamoDB event.", err)
		return
	}
	newBattle := battles.Battle{}
	err = dynamodbattribute.UnmarshalMap(record.Change.NewImage, &newBattle)
	if err != nil {
		handler.logger.Error("Failed to unmarshal new image in DynamoDB event.", err)
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
	var err error
	var score float64
	if change.createdOnChanged {
		score = handler.calculateRecency(change.newCreatedOn)
		err = handler.index.Set(battles.RecentSort, battleID, score)
		if err != nil {
			handler.logger.Error("Failed to set value in "+battles.RecentSort+" index.", err)
		}
	}
	if change.votesChanged {
		score = handler.calculatePopularity(change.newVotesFor, change.newVotesAgainst)
		err = handler.index.Set(battles.PopularSort, battleID, score)
		if err != nil {
			handler.logger.Error("Failed to set value in "+battles.PopularSort+" index.", err)
		}
		score = handler.calculateControversy(change.newVotesFor, change.newVotesAgainst)
		err = handler.index.Set(battles.ControversialSort, battleID, score)
		if err != nil {
			handler.logger.Error("Failed to set value in "+battles.ControversialSort+" index.", err)
		}
	}
}

func (handler *Handler) calculateRecency(createdOn int64) float64 {
	return float64(createdOn)
}

func (handler *Handler) calculatePopularity(votesFor int, votesAgainst int) float64 {
	totalVotes := float64(votesFor + votesAgainst)
	return handler.getRecencyWeight() + totalVotes
}

func (handler *Handler) calculateControversy(votesFor int, votesAgainst int) float64 {
	totalVotes := float64(votesFor + votesAgainst)
	voteDifference := math.Abs(float64(votesFor - votesAgainst))
	return handler.getRecencyWeight() + totalVotes - voteDifference
}

func (handler *Handler) getRecencyWeight() float64 {
	return float64(time.Now().Unix() / 86400)
}
