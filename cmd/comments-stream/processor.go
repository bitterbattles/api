package main

import (
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/comments"
	"github.com/bitterbattles/api/pkg/lambda/stream"
)

// Processor represents a stream event processor
type Processor struct {
	commentsIndex     comments.IndexInterface
	battlesRepository battles.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(commentsIndex comments.IndexInterface, battlesRepository battles.RepositoryInterface) *Processor {
	return &Processor{
		commentsIndex:     commentsIndex,
		battlesRepository: battlesRepository,
	}
}

// Process processes a DynamoDB event
func (processor *Processor) Process(event *stream.Event) error {
	newBattleComments := make(map[string][]*comments.Comment)
	for _, record := range event.Records {
		processor.captureChanges(&record, newBattleComments)
	}
	for battleID, newComments := range newBattleComments {
		processor.processNewComments(battleID, newComments)
	}
	return nil
}

func (processor *Processor) captureChanges(record *stream.EventRecord, newBattleComments map[string][]*comments.Comment) {
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
	commentID := newComment.ID
	if commentID == "" {
		log.Println("Unexpected missing new comment ID.")
		return
	}
	if oldComment.ID == "" {
		// New comment
		battleID := newComment.BattleID
		if battleID == "" {
			log.Println("Unexpected missing new comment Battle ID.")
			return
		}
		newComments := newBattleComments[battleID]
		if newComments == nil {
			newComments = make([]*comments.Comment, 0, 1)
		}
		newComments = append(newComments, newComment)
		newBattleComments[battleID] = newComments
	}
}

func (processor *Processor) processNewComments(battleID string, newComments []*comments.Comment) {
	for _, comment := range newComments {
		err := processor.commentsIndex.Upsert(comment)
		if err != nil {
			log.Println("Failed to add new comment to indexes. Error:", err)
		}
	}
	err := processor.battlesRepository.IncrementComments(battleID, len(newComments))
	if err != nil {
		log.Println("Failed to increment Battle comments. Error:", err)
	}
}
