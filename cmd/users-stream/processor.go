package main

import (
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/comments"
	"github.com/bitterbattles/api/pkg/lambda/stream"
	"github.com/bitterbattles/api/pkg/users"
)

// Processor represents a stream processor
type Processor struct {
	battlesRepository  battles.RepositoryInterface
	commentsRepository comments.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(battlesRepository battles.RepositoryInterface, commentsRepository comments.RepositoryInterface) *Processor {
	return &Processor{
		battlesRepository:  battlesRepository,
		commentsRepository: commentsRepository,
	}
}

// Process processes a DynamoDB event
func (processor *Processor) Process(event *stream.Event) error {
	changes := make(map[string]*change)
	for _, record := range event.Records {
		processor.captureChange(&record, changes)
	}
	for userID, change := range changes {
		processor.processChange(userID, change)
	}
	return nil
}

func (processor *Processor) captureChange(record *stream.EventRecord, changes map[string]*change) {
	var err error
	oldUser := &users.User{}
	err = dynamodbattribute.UnmarshalMap(record.Change.OldImage, oldUser)
	if err != nil {
		log.Println("Failed to unmarshal old image in DynamoDB event. Error:", err)
		return
	}
	newUser := &users.User{}
	err = dynamodbattribute.UnmarshalMap(record.Change.NewImage, newUser)
	if err != nil {
		log.Println("Failed to unmarshal new image in DynamoDB event. Error:", err)
		return
	}
	id := newUser.ID
	if id == "" {
		id = oldUser.ID
		if id == "" {
			log.Println("Unexpected missing new User ID.")
			return
		}
	}
	c := changes[id]
	if c == nil {
		c = &change{
			oldUser: oldUser,
			newUser: newUser,
		}
	} else {
		c.newUser = newUser
	}
	changes[id] = c
}

func (processor *Processor) processChange(userID string, change *change) {
	var err error
	newUser := change.newUser
	if newUser.ID == "" {
		// Deleted user
		err = processor.battlesRepository.UpdateUsername(userID, users.AnonymousUsername)
		if err != nil {
			log.Println("Failed to update username across battles. Error:", err)
		}
		err = processor.commentsRepository.UpdateUsername(userID, users.AnonymousUsername)
		if err != nil {
			log.Println("Failed to update username across comments. Error:", err)
		}
	}
}
