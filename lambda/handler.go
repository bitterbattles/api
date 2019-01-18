package lambda

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/core/errors"
	"github.com/bitterbattles/api/core/loggers"
)

// Handler represents a lambda handler
type Handler struct {
	function interface{}
	logger   loggers.LoggerInterface
}

// NewHandler creates a new handler instance
func NewHandler(function interface{}, logger loggers.LoggerInterface) *Handler {
	return &Handler{function, logger}
}

// Invoke processes the request and creates the response
func (handler Handler) Invoke(context context.Context, requestBytes []byte) (responseBytes []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			responseBytes, err = handler.internalError("Unexpected panic.", err)
		}
	}()
	if handler.function == nil {
		return handler.internalError("Missing function.", nil)
	}
	functionType := reflect.TypeOf(handler.function)
	if functionType.Kind() != reflect.Func {
		return handler.internalError("Unexpected function type.", nil)
	}
	if functionType.NumIn() != 1 {
		return handler.internalError("Unexpected number of arguments in function.", nil)
	}
	request := reflect.New(functionType.In(0))
	err = json.Unmarshal(requestBytes, request.Interface())
	if err != nil {
		return handler.internalError("Failed to decode request JSON.", err)
	}
	args := []reflect.Value{request.Elem()}
	returnArgs := reflect.ValueOf(handler.function).Call(args)
	numReturnArgs := len(returnArgs)
	if !(numReturnArgs == 1 || numReturnArgs == 2) {
		return handler.internalError("Unexpected number of return values from handler.", nil)
	}
	var returnArg = returnArgs[numReturnArgs-1].Interface()
	var response interface{}
	if returnArg != nil {
		returnError, ok := returnArg.(error)
		if !ok {
			return handler.internalError("Unexpected non-error return value.", nil)
		}
		httpError, ok := returnArg.(errors.HTTPError)
		if !ok {
			return handler.internalError("Handler returned an unexpected error.", returnError)
		}
		response = httpError
	} else if numReturnArgs == 2 {
		response = returnArgs[0].Interface()
	}
	if response != nil {
		responseBytes, err = json.Marshal(response)
		if err != nil {
			return handler.internalError("Failed to encode response JSON.", err)
		}
	}
	return responseBytes, nil
}

// Run runs the lambda handler
func (handler Handler) Run() {
	lambda.StartHandler(handler)
}

func (handler Handler) internalError(message string, err error) ([]byte, error) {
	handler.logger.Error(message, err)
	response := errors.NewInternalServerError()
	return json.Marshal(response)
}
