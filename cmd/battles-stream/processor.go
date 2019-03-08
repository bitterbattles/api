package main

import (
	"log"
	"math"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/lambda/stream"
	"github.com/bitterbattles/api/pkg/ranks"
	"github.com/bitterbattles/api/pkg/time"
)

// Processor represents a stream event processor
type Processor struct {
	repository ranks.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(repository ranks.RepositoryInterface) *Processor {
	return &Processor{
		repository: repository,
	}
}

// Process processes a DynamoDB event
func (processor *Processor) Process(event *stream.Event) error {
	changes := make(map[string]*change)
	for _, record := range event.Records {
		processor.captureChange(&record, changes)
	}
	for battleID, change := range changes {
		processor.processChange(battleID, change)
	}
	return nil
}

func (processor *Processor) captureChange(record *stream.EventRecord, changes map[string]*change) {
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

func (processor *Processor) processChange(battleID string, change *change) {
	var err error
	var score float64
	if change.deleted {
		processor.repository.DeleteByBattleID(battles.RecentSort, battleID)
		processor.repository.DeleteByBattleID(battles.PopularSort, battleID)
		processor.repository.DeleteByBattleID(battles.ControversialSort, battleID)
	} else {
		if change.createdOnChanged {
			score = processor.calculateRecency(change.newCreatedOn)
			err = processor.repository.SetScore(battles.RecentSort, battleID, score)
			if err != nil {
				log.Println("Failed to set value in", battles.RecentSort, "ranking. Error:", err)
			}
		}
		if change.votesChanged {
			score = processor.calculatePopularity(change.newVotesFor, change.newVotesAgainst)
			err = processor.repository.SetScore(battles.PopularSort, battleID, score)
			if err != nil {
				log.Println("Failed to set value in", battles.PopularSort, "ranking. Error:", err)
			}
			score = processor.calculateControversy(change.newVotesFor, change.newVotesAgainst)
			err = processor.repository.SetScore(battles.ControversialSort, battleID, score)
			if err != nil {
				log.Println("Failed to set value in", battles.ControversialSort, "ranking. Error:", err)
			}
		}
	}
}

func (processor *Processor) calculateRecency(createdOn int64) float64 {
	return float64(createdOn)
}

func (processor *Processor) calculatePopularity(votesFor int, votesAgainst int) float64 {
	totalVotes := float64(votesFor + votesAgainst)
	return processor.getRecencyWeight() + totalVotes
}

func (processor *Processor) calculateControversy(votesFor int, votesAgainst int) float64 {
	totalVotes := float64(votesFor + votesAgainst)
	voteDifference := math.Abs(float64(votesFor - votesAgainst))
	return processor.getRecencyWeight() + totalVotes - voteDifference
}

func (processor *Processor) getRecencyWeight() float64 {
	return float64(time.NowUnix() / 86400)
}
