package votes_test

import (
	"os"
	"testing"

	"github.com/bitterbattles/api/core/errors"
	"github.com/bitterbattles/api/core/tests"
	"github.com/bitterbattles/api/votes"
	"github.com/bitterbattles/api/votes/mocks"
)

var repository *mocks.Repository
var manager *votes.Manager
var controller *votes.Controller

func TestMain(m *testing.M) {
	repository = mocks.NewRepository()
	manager = votes.NewManager(repository)
	controller = votes.NewController(manager)
	os.Exit(m.Run())
}

func TestPost(t *testing.T) {
	reset()
	battleID := "battleIdbattleIdbatt"
	isVoteFor := true
	err := post(battleID, isVoteFor)
	tests.AssertNil(t, err)
	vote := repository.GetLastAdded()
	tests.AssertNotNil(t, vote)
	tests.AssertEquals(t, vote.BattleID, battleID)
	tests.AssertEquals(t, vote.IsVoteFor, isVoteFor)
}

func TestPostNoBattleId(t *testing.T) {
	reset()
	err := post("", true)
	tests.AssertHTTPError(t, err, errors.BadRequestCode)
	tests.AssertNil(t, repository.GetLastAdded())
}

func reset() {
	repository.Reset()
}

func post(battleID string, isVoteFor bool) error {
	request := votes.PostRequest{
		BattleID:  battleID,
		IsVoteFor: isVoteFor,
	}
	return controller.HandlePost(request)
}
