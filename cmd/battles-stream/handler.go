package main

import (
	"log"
	"math"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/common/lambda/stream"
	"github.com/bitterbattles/api/pkg/ranks"
)

// Handler represents a stream handler
type Handler struct {
	repository ranks.RepositoryInterface
}

// NewHandler creates a new Handler instance
func NewHandler(repository ranks.RepositoryInterface) *stream.Handler {
	handler := Handler{
		repository: repository,
	}
	return stream.NewHandler(&handler)
}

// Handle handles a DynamoDB event
func (handler *Handler) Handle(event *stream.Event) error {
	changes := make(map[string]*change)
	for _, record := range event.Records {
		handler.captureChange(&record, changes)
	}
	for battleID, change := range changes {
		handler.processChange(battleID, change)
	}
	return nil
}

func (handler *Handler) captureChange(record *stream.EventRecord, changes map[string]*change) {
	var err error
	oldBattle := battles.Battle{}
	err = dynamodbattribute.UnmarshalMap(record.Change.OldImage, &oldBattle)
	if err != nil {
		log.Println("Failed to unmarshal old image in DynamoDB event. Error:", err)
		return
	}
	newBattle := battles.Battle{}
	err = dynamodbattribute.UnmarshalMap(record.Change.NewImage, &newBattle)
	if err != nil {
		log.Println("Failed to unmarshal new image in DynamoDB event. Error:", err)
		return
	}
	id := newBattle.ID
	if id == "" {
		log.Println("Unexpected missing new Battle ID.")
		return
	}
	c := changes[id]
	if c == nil {
		c = &change{}
	}
	if newBattle.State == battles.Deleted {
		if oldBattle.State == battles.Deleted {
			log.Println("Unexpected modification of deleted battle.")
			return
		}
		c.deleted = true
	} else {
		if newBattle.CreatedOn != oldBattle.CreatedOn {
			c.createdOnChanged = true
			c.newCreatedOn = newBattle.CreatedOn
		}
		if oldBattle.ID == "" || newBattle.VotesFor != oldBattle.VotesFor || newBattle.VotesAgainst != oldBattle.VotesAgainst {
			c.votesChanged = true
			c.newVotesFor = newBattle.VotesFor
			c.newVotesAgainst = newBattle.VotesAgainst
		}
	}
	changes[id] = c
}

func (handler *Handler) processChange(battleID string, change *change) {
	var err error
	var score float64
	if change.deleted {
		handler.repository.DeleteByBattleID(battles.RecentSort, battleID)
		handler.repository.DeleteByBattleID(battles.PopularSort, battleID)
		handler.repository.DeleteByBattleID(battles.ControversialSort, battleID)
	} else {
		if change.createdOnChanged {
			score = handler.calculateRecency(change.newCreatedOn)
			err = handler.repository.SetScore(battles.RecentSort, battleID, score)
			if err != nil {
				log.Println("Failed to set value in", battles.RecentSort, "ranking. Error:", err)
			}
		}
		if change.votesChanged {
			score = handler.calculatePopularity(change.newVotesFor, change.newVotesAgainst)
			err = handler.repository.SetScore(battles.PopularSort, battleID, score)
			if err != nil {
				log.Println("Failed to set value in", battles.PopularSort, "ranking. Error:", err)
			}
			score = handler.calculateControversy(change.newVotesFor, change.newVotesAgainst)
			err = handler.repository.SetScore(battles.ControversialSort, battleID, score)
			if err != nil {
				log.Println("Failed to set value in", battles.ControversialSort, "ranking. Error:", err)
			}
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
