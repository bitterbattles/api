package main

import (
	"github.com/bitterbattles/api/pkg/errors"
	"github.com/bitterbattles/api/pkg/jwt"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/loginspost"
	"github.com/bitterbattles/api/pkg/time"
	"github.com/bitterbattles/api/pkg/users"
)

// Processor represents a request processor
type Processor struct {
	repository         users.RepositoryInterface
	accessTokenSecret  string
	refreshTokenSecret string
}

// NewProcessor creates a new Processor instance
func NewProcessor(repository users.RepositoryInterface, accessTokenSecret string, refreshTokenSecret string) *Processor {
	return &Processor{
		repository:         repository,
		accessTokenSecret:  accessTokenSecret,
		refreshTokenSecret: refreshTokenSecret,
	}
}

// NewRequestBody creates a new request body instance
func (processor *Processor) NewRequestBody() interface{} {
	return &Request{}
}

// Process processes a request
func (processor *Processor) Process(input *api.Input) (*api.Output, error) {
	request, _ := input.RequestBody.(*Request)
	authContext := &api.AuthContext{}
	err := jwt.DecodeHS256(request.RefreshToken, processor.refreshTokenSecret, authContext)
	if err != nil {
		return nil, errors.NewBadRequestError("Invalid token.")
	}
	if authContext.ExpiresOn <= time.NowUnix() {
		return nil, errors.NewBadRequestError("Expired token.")
	}
	user, err := processor.repository.GetByID(authContext.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NewBadRequestError("User not found.")
	}
	response, err := loginspost.CreateResponse(user.ID, processor.accessTokenSecret, processor.refreshTokenSecret)
	if err != nil {
		return nil, err
	}
	output := api.NewOutput(response)
	return output, nil
}
