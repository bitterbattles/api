package votesstream

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bitterbattles/api/battles"
	"github.com/bitterbattles/api/common/handlers"
	"github.com/bitterbattles/api/votes"
)

// Handler represents a stream handler
type Handler struct {
	repository battles.RepositoryInterface
}

// NewHandler creates a new Handler instance
func NewHandler(repository battles.RepositoryInterface) *handlers.StreamHandler {
	handler := Handler{
		repository: repository,
	}
	return handlers.NewStreamHandler(&handler)
}

// TODO: Log errors

// Handle handles a DynamoDB event
func (handler *Handler) Handle(event *handlers.DynamoEvent) error {
	changes := make(map[string]*change)
	for _, record := range event.Records {
		handler.captureChange(&record, changes)
	}
	for battleID, change := range changes {
		handler.repository.IncrementVotes(battleID, change.deltaVotesFor, change.deltaVotesAgainst)
	}
	return nil
}

func (handler *Handler) captureChange(record *handlers.DynamoEventRecord, changes map[string]*change) {
	var err error
	vote := votes.Vote{}
	err = dynamodbattribute.UnmarshalMap(record.Change.NewImage, &vote)
	if err != nil {
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
