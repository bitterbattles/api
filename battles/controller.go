package battles

import (
	"github.com/bitterbattles/api/battles/errors"
	coreErrors "github.com/bitterbattles/api/core/errors"
)

// GetRequest represents a GET request
type GetRequest struct {
	Sort     string `json:"sort"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

// GetResponse represents a GET response
type GetResponse struct {
	ID           string `json:"id"`
	UserID       string `json:"-"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	VotesFor     int    `json:"votesFor"`
	VotesAgainst int    `json:"votesAgainst"`
	CreatedOn    int64  `json:"createdOn"`
}

// PostRequest represents a POST request
type PostRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Controller is used to handle requests against the "battles" resource
type Controller struct {
	manager ManagerInterface
}

// NewController creates a new Controller instance
func NewController(manager ManagerInterface) *Controller {
	return &Controller{manager}
}

// HandleGet handles GET requests
func (controller *Controller) HandleGet(request GetRequest) ([]GetResponse, error) {
	battles, err := controller.manager.GetPage(request.Sort, request.Page, request.PageSize)
	if err != nil {
		return nil, err
	}
	count := len(battles)
	response := make([]GetResponse, count)
	for i := 0; i < count; i++ {
		response[i] = GetResponse(*battles[i])
	}
	return response, nil
}

// HandlePost handles POST requests
func (controller *Controller) HandlePost(request PostRequest) error {
	err := controller.manager.Create(request.Title, request.Description)
	if err != nil {
		if _, ok := err.(errors.InvalidTitleError); ok {
			return coreErrors.NewBadRequestError(err.Error())
		} else if _, ok := err.(errors.InvalidDescriptionError); ok {
			return coreErrors.NewBadRequestError(err.Error())
		}
		return err
	}
	return nil
}
