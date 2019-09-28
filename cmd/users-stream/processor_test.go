package main_test

import (
	"encoding/json"
	"testing"

	. "github.com/bitterbattles/api/cmd/users-stream"
	"github.com/bitterbattles/api/pkg/battles"
	battlesMocks "github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/comments"
	commentsMocks "github.com/bitterbattles/api/pkg/comments/mocks"
	"github.com/bitterbattles/api/pkg/lambda/stream"
	. "github.com/bitterbattles/api/pkg/tests"
)

func TestProcessor(t *testing.T) {
	battlesRepository := battlesMocks.NewRepository()
	battle := &battles.Battle{
		UserID:   "userId0",
		Username: "username0",
	}
	battlesRepository.Add(battle)
	commentsRepository := commentsMocks.NewRepository()
	comment := &comments.Comment{
		UserID:   "userId0",
		Username: "username0",
	}
	commentsRepository.Add(comment)
	processor := NewProcessor(battlesRepository, commentsRepository)
	eventJSON := `{"records":[{"dynamodb":{"OldImage":{"id":{"S":"userId0"}}}}]}`
	event := newEvent(eventJSON)
	err := processor.Process(event)
	AssertNil(t, err)
	foundBattle := battlesRepository.GetLastAdded()
	AssertEquals(t, foundBattle.Username, "[Deleted]")
	foundComment := commentsRepository.GetLastAdded()
	AssertEquals(t, foundComment.Username, "[Deleted]")
}

func newEvent(eventJSON string) *stream.Event {
	event := &stream.Event{}
	json.Unmarshal([]byte(eventJSON), event)
	return event
}
