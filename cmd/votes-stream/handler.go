package main

import (
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/common/lambda/stream"
	"github.com/bitterbattles/api/pkg/votes"
)

// Handler represents a stream handler
type Handler struct {
	repository battles.RepositoryInterface
}

// NewHandler creates a new Handler instance
func NewHandler(repository battles.RepositoryInterface) *stream.Handler {
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
		err := handler.repository.IncrementVotes(battleID, change.deltaVotesFor, change.deltaVotesAgainst)
		if err != nil {
			log.Println("Failed to increment votes for battle ID", battleID, ".")
		}
	}
	return nil
}

func (handler *Handler) captureChange(record *stream.EventRecord, changes map[string]*change) {
	var err error
	vote := votes.Vote{}
	err = dynamodbattribute.UnmarshalMap(record.Change.NewImage, &vote)
	if err != nil {
		log.Println("Failed to unmarshal new image in DynamoDB event. Error:", err)
		return
	}
	battleID := vote.BattleID
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
