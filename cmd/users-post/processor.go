package main

import (
	"strings"

	"github.com/bitterbattles/api/pkg/crypto"
	"github.com/bitterbattles/api/pkg/errors"
	"github.com/bitterbattles/api/pkg/guid"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/input"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/time"
	"github.com/bitterbattles/api/pkg/users"
)

const (
	minUsernameLength  = 4
	maxUsernameLength  = 20
	usernameRegex      = `^[a-zA-Z][a-zA-Z0-9]*$`
	minPasswordLength  = 8
	maxPasswordLength  = 24
	minPasswordUpper   = 1
	minPasswordNumbers = 1
	minPasswordSymbols = 1
)

// Processor represents a request processor
type Processor struct {
	repository users.RepositoryInterface
}

// NewProcessor creates a new Processor instance
func NewProcessor(repository users.RepositoryInterface) *Processor {
	return &Processor{
		repository: repository,
	}
}

// NewRequestBody creates a new request body instance
func (processor *Processor) NewRequestBody() interface{} {
	return &Request{}
}

// Process processes a request
func (processor *Processor) Process(input *api.Input) (*api.Output, error) {
	request, _ := input.RequestBody.(*Request)
	username, err := processor.sanitizeUsername(request.Username)
	if err != nil {
		return nil, err
	}
	password, err := processor.sanitizePassword(request.Password)
	if err != nil {
		return nil, err
	}
	existingUser, err := processor.repository.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.NewBadRequestErrorWithCode(errors.UsernameAlreadyExists, "The username is already taken.")
	}
	passwordHash, err := crypto.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := users.User{
		ID:              guid.New(),
		Username:        strings.ToLower(username),
		DisplayUsername: username,
		PasswordHash:    passwordHash,
		CreatedOn:       time.NowUnix(),
		State:           users.Active,
	}
	err = processor.repository.Add(&user)
	if err != nil {
		return nil, err
	}
	output := api.NewOutputWithStatus(http.Created, nil)
	return output, nil
}

func (processor *Processor) sanitizeUsername(value string) (string, error) {
	rules := input.StringRules{
		MinLength: minUsernameLength,
		MaxLength: maxUsernameLength,
		Regex:     usernameRegex,
	}
	errorCreator := func(message string) error {
		return errors.NewBadRequestError("Invalid username: " + message)
	}
	return input.SanitizeString(value, rules, errorCreator)
}

func (processor *Processor) sanitizePassword(value string) (string, error) {
	rules := input.PasswordRules{
		MinLength:  minPasswordLength,
		MaxLength:  maxPasswordLength,
		MinUpper:   minPasswordUpper,
		MinNumbers: minPasswordNumbers,
		MinSymbols: minPasswordSymbols,
	}
	errorCreator := func(message string) error {
		return errors.NewBadRequestError("Invalid password: " + message)
	}
	return input.SanitizePassword(value, rules, errorCreator)
}
