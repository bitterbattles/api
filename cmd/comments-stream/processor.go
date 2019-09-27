package main

import (
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bitterbattles/api/pkg/comments"
	"github.com/bitterbattles/api/pkg/lambda/stream"
)

// Processor represents a stream event processor
type Processor struct {
	indexer *comments.Indexer
}

// NewProcessor creates a new Processor instance
func NewProcessor(indexer *comments.Indexer) *Processor {
	return &Processor{
		indexer: indexer,
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
	oldComment := &comments.Comment{}
	err = dynamodbattribute.UnmarshalMap(record.Change.OldImage, oldComment)
	if err != nil {
		log.Println("Failed to unmarshal old image in DynamoDB event. Error:", err)
		return
	}
	newComment := &comments.Comment{}
	err = dynamodbattribute.UnmarshalMap(record.Change.NewImage, newComment)
	if err != nil {
		log.Println("Failed to unmarshal new image in DynamoDB event. Error:", err)
		return
	}
	id := newComment.ID
	if id == "" {
		log.Println("Unexpected missing new Comment ID.")
		return
	}
	c := changes[id]
	if c == nil {
		c = &change{
			oldComment: oldComment,
			newComment: newComment,
		}
	} else {
		c.newComment = newComment
		changes[id] = c
	}
	changes[id] = c
}

func (processor *Processor) processChange(change *change) {
	var err error
	oldComment := change.oldComment
	newComment := change.newComment
	if oldComment.State == comments.Deleted {
		log.Println("Unexpected modification of deleted comment ID", oldComment.ID)
		return
	}
	if newComment.State != oldComment.State && newComment.State == comments.Deleted {
		// Deleted comment
		err = processor.indexer.Delete(newComment)
		if err != nil {
			log.Println("Failed to delete comment from indexes. Error:", err)
		}
	} else if oldComment.ID == "" {
		// New comment
		err = processor.indexer.Add(newComment)
		if err != nil {
			log.Println("Failed to add new comment to indexes. Error:", err)
		}
	}
}
