package votes

import (
	coreErrors "github.com/bitterbattles/api/core/errors"
	"github.com/bitterbattles/api/votes/errors"
)

// PostRequest represents a POST request
type PostRequest struct {
	BattleID  string `json:"battleId"`
	IsVoteFor bool   `json:"isVoteFor"`
}

// Controller is used to handle requests against the "votes" resource
type Controller struct {
	manager ManagerInterface
}

// NewController creates a new Controller instance
func NewController(manager ManagerInterface) *Controller {
	return &Controller{manager}
}

// HandlePost handles POST requests
func (controller *Controller) HandlePost(request PostRequest) error {
	err := controller.manager.Create(request.BattleID, request.IsVoteFor)
	if err != nil {
		if _, ok := err.(errors.InvalidBattleIDError); ok {
			return coreErrors.NewBadRequestError(err.Error())
		}
		return err
	}
	return nil
}
