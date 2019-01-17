package lambda

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/core/errors"
)

// Handler represents a lambda handler
type Handler struct {
	function interface{}
}

// NewHandler creates a new handler instance
func NewHandler(function interface{}) *Handler {
	return &Handler{function}
}

// Invoke processes the request and creates the response
func (handler Handler) Invoke(context context.Context, requestBytes []byte) ([]byte, error) {
	// TODO: Catch panics with defer/recover
	if handler.function == nil {
		return internalError("Missing function.", nil)
	}
	functionType := reflect.TypeOf(handler.function)
	if functionType.Kind() != reflect.Func {
		return internalError("Unexpected function type.", nil)
	}
	if functionType.NumIn() != 1 {
		return internalError("Unexpected number of arguments in function.", nil)
	}
	request := reflect.New(functionType.In(0))
	err := json.Unmarshal(requestBytes, request.Interface())
	if err != nil {
		return internalError("Failed to decode request JSON.", err)
	}
	args := []reflect.Value{request.Elem()}
	callReturn := reflect.ValueOf(handler.function).Call(args)
	if len(callReturn) != 2 {
		return internalError("Unexpected number of return values from handler.", nil)
	}
	var callError = callReturn[1].Interface()
	var response interface{}
	if callError != nil {
		httpError, ok := callError.(errors.HTTPError)
		if !ok {
			return internalError("Handler returned an unexpected error.", err)
		}
		response = httpError
	} else {
		response = callReturn[0].Interface()
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return internalError("Failed to encode response JSON.", err)
	}
	return responseBytes, nil
}

// Run runs the lambda handler
func (handler Handler) Run() {
	lambda.StartHandler(handler)
}

func internalError(message string, err error) ([]byte, error) {
	// TODO: Log message and error
	response := errors.NewInternalServerError()
	return json.Marshal(response)
}
