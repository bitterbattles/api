package main_test

import (
	"encoding/json"
	"testing"

	. "github.com/bitterbattles/api/cmd/comments-stream"
	"github.com/bitterbattles/api/pkg/battles"
	battlesMocks "github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/comments"
	commentsMocks "github.com/bitterbattles/api/pkg/comments/mocks"
	"github.com/bitterbattles/api/pkg/lambda/stream"
	. "github.com/bitterbattles/api/pkg/tests"
)

func TestProcessorNewComment(t *testing.T) {
	commentsIndex := commentsMocks.NewIndex()
	battlesRepository := battlesMocks.NewRepository()
	battle := &battles.Battle{
		ID:       "battleId0",
		Comments: 0,
	}
	battlesRepository.Add(battle)
	processor := NewProcessor(commentsIndex, battlesRepository)
	eventJSON := `{"records":[{"dynamodb":{"NewImage":{"id":{"S":"id0"},"battleId":{"S":"battleId0"},"userId":{"S":"userId0"},"comment":{"S":"comment0"},"createdOn":{"N":"123"}}}}]}`
	event := newEvent(eventJSON)
	err := processor.Process(event)
	AssertNil(t, err)
	expectedComment := &comments.Comment{
		ID:        "id0",
		BattleID:  "battleId0",
		UserID:    "userId0",
		Comment:   "comment0",
		CreatedOn: 123,
	}
	verifyIndex(t, commentsIndex, expectedComment)
	verifyBattleComments(t, battlesRepository, "battleId0", 1)
}

func newEvent(eventJSON string) *stream.Event {
	event := &stream.Event{}
	json.Unmarshal([]byte(eventJSON), event)
	return event
}

func verifyIndex(t *testing.T, index *commentsMocks.Index, expectedComment *comments.Comment) {
	comment := index.GetLastCommentAdded()
	if expectedComment == nil {
		AssertNil(t, comment)
	} else {
		AssertNotNil(t, comment)
		AssertEquals(t, comment.ID, expectedComment.ID)
		AssertEquals(t, comment.BattleID, expectedComment.BattleID)
		AssertEquals(t, comment.UserID, expectedComment.UserID)
		AssertEquals(t, comment.Comment, expectedComment.Comment)
	}
}

func verifyBattleComments(t *testing.T, repository *battlesMocks.Repository, id string, expectedComments int) {
	battle, err := repository.GetByID(id)
	AssertNil(t, err)
	AssertNotNil(t, battle)
	AssertEquals(t, battle.Comments, expectedComments)
}
