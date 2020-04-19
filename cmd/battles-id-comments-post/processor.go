package main

import (
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/comments"
	"github.com/bitterbattles/api/pkg/errors"
	"github.com/bitterbattles/api/pkg/guid"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/input"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/time"
	"github.com/bitterbattles/api/pkg/users"
)

const (
	idParam          = "id"
	minCommentLength = 4
	maxCommentLength = 500
)

// Processor represents a request processor
type Processor struct {
	commentsRepository comments.RepositoryInterface
	usersRepository    users.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(commentsRepository comments.RepositoryInterface, usersRepository users.RepositoryInterface) *Processor {
	return &Processor{
		commentsRepository: commentsRepository,
		usersRepository:    usersRepository,
	}
}

// NewRequestBody creates a new request body instance
func (processor *Processor) NewRequestBody() interface{} {
	return &Request{}
}

// Process processes a request
func (processor *Processor) Process(input *api.Input) (*api.Output, error) {
	request, _ := input.RequestBody.(*Request)
	battleID := input.PathParams[idParam]
	userID := input.AuthContext.UserID
	user, err := processor.usersRepository.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NewForbiddenError("User not found.")
	}
	commentText, err := sanitizeComment(request.Comment)
	if err != nil {
		return nil, err
	}
	comment := &comments.Comment{
		ID:        guid.New(),
		BattleID:  battleID,
		UserID:    userID,
		Comment:   commentText,
		CreatedOn: time.NowUnix(),
		State:     battles.Active,
	}
	err = processor.commentsRepository.Add(comment)
	if err != nil {
		return nil, err
	}
	output := api.NewOutputWithStatus(http.Created, nil)
	return output, nil
}

func sanitizeComment(comment string) (string, error) {
	rules := input.StringRules{
		TrimSpace: true,
		MinLength: minCommentLength,
		MaxLength: maxCommentLength,
	}
	errorCreator := func(message string) error {
		return errors.NewBadRequestError("Invalid comment: " + message)
	}
	return input.SanitizeString(comment, rules, errorCreator)
}
