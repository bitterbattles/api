package main

import (
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/lambda/stream"
	"github.com/bitterbattles/api/pkg/votes"
)

// Processor represents a stream processor
type Processor struct {
	repository battles.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(repository battles.RepositoryInterface) *Processor {
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
		err := processor.repository.IncrementVotes(battleID, change.deltaVotesFor, change.deltaVotesAgainst)
		if err != nil {
			log.Println("Failed to increment votes. Error:", err)
		}
	}
	return nil
}

func (processor *Processor) captureChange(record *stream.EventRecord, changes map[string]*change) {
	var err error
	vote := votes.Vote{}
	err = dynamodbattribute.UnmarshalMap(record.Change.NewImage, &vote)
	if err != nil {
		log.Println("Failed to unmarshal new image in DynamoDB event. Error:", err)
		return
	}
	battleID := vote.BattleID
	if battleID == "" {
		log.Println("Unexpected missing battle ID.")
		return
	}
	c, ok := changes[battleID]
	if !ok {
		c = &change{}
	}
	if vote.IsVoteFor {
		c.deltaVotesFor++
	} else {
		c.deltaVotesAgainst++
	}
	changes[battleID] = c
}
