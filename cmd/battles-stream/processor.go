package main

import (
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/lambda/stream"
)

// Processor represents a stream event processor
type Processor struct {
	battlesIndex  battles.IndexInterface
	battlesScorer *battles.Scorer
}

// NewProcessor creates a new Processor instance
func NewProcessor(battlesIndex battles.IndexInterface, battlesScorer *battles.Scorer) *Processor {
	return &Processor{
		battlesIndex:  battlesIndex,
		battlesScorer: battlesScorer,
	}
}

// Process processes a DynamoDB event
func (processor *Processor) Process(event *stream.Event) error {
	changes := make(map[string]*change)
	for _, record := range event.Records {
		processor.captureChange(&record, changes)
	}
	for _, change := range changes {
		processor.processChange(change)
	}
	return nil
}

func (processor *Processor) captureChange(record *stream.EventRecord, changes map[string]*change) {
	var err error
	oldBattle := &battles.Battle{}
	err = dynamodbattribute.UnmarshalMap(record.Change.OldImage, oldBattle)
	if err != nil {
		log.Println("Failed to unmarshal old image in DynamoDB event. Error:", err)
		return
	}
	newBattle := &battles.Battle{}
	err = dynamodbattribute.UnmarshalMap(record.Change.NewImage, newBattle)
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
		c = &change{
			oldBattle: oldBattle,
			newBattle: newBattle,
		}
	} else {
		c.newBattle = newBattle
		changes[id] = c
	}
	changes[id] = c
}

func (processor *Processor) processChange(change *change) {
	var err error
	oldBattle := change.oldBattle
	newBattle := change.newBattle
	if oldBattle.State == battles.Deleted {
		log.Println("Unexpected modification of deleted battle ID", oldBattle.ID)
		return
	}
	if newBattle.State != oldBattle.State && newBattle.State == battles.Deleted {
		// Deleted battle
		err = processor.battlesIndex.Delete(newBattle)
		if err != nil {
			log.Println("Failed to delete battle from index. Error:", err)
		}
	} else {
		// New or updated battle
		popularity := processor.battlesScorer.ScorePopularity(newBattle)
		controversy := processor.battlesScorer.ScoreControversy(newBattle)
		err = processor.battlesIndex.Upsert(newBattle, popularity, controversy)
		if err != nil {
			log.Println("Failed to index battle. Error:", err)
		}
	}
}
