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
	indexer    *comments.Indexer
	repository battles.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(indexer *comments.Indexer, repository battles.RepositoryInterface) *Processor {
	return &Processor{
		indexer:    indexer,
		repository: repository,
	}
}

// Process processes a DynamoDB event
func (processor *Processor) Process(event *stream.Event) error {
	changedComments := make(map[string]*changedComment)
	changedBattles := make(map[string]*changedBattle)
	for _, record := range event.Records {
		processor.captureChanges(&record, changedComments, changedBattles)
	}
	for _, change := range changedComments {
		processor.processCommentChange(change)
	}
	for battleID, change := range changedBattles {
		processor.processBattleChange(battleID, change)
	}
	return nil
}

func (processor *Processor) captureChanges(record *stream.EventRecord, changedComments map[string]*changedComment, changedBattles map[string]*changedBattle) {
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
		log.Println("Unexpected missing new Comment ID.")
		return
	}
	commentChange := changedComments[commentID]
	if commentChange == nil {
		commentChange = &changedComment{
			oldComment: oldComment,
			newComment: newComment,
		}
	} else {
		commentChange.newComment = newComment
		changedComments[commentID] = commentChange
	}
	changedComments[commentID] = commentChange
	battleID := newComment.BattleID
	if battleID == "" {
		log.Println("Unexpected missing new Battle ID.")
		return
	}
	battleChange := changedBattles[battleID]
	if battleChange == nil {
		battleChange = &changedBattle{}
	}
	battleChange.deltaComments++
	changedBattles[battleID] = battleChange
}

func (processor *Processor) processCommentChange(change *changedComment) {
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

func (processor *Processor) processBattleChange(battleID string, change *changedBattle) {
	var err error
	err = processor.repository.IncrementComments(battleID, change.deltaComments)
	if err != nil {
		log.Println("Failed to increment comments. Error:", err)
	}
}
