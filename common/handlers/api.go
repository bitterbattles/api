package handlers

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/bitterbattles/api/common/errors"
	"github.com/bitterbattles/api/common/loggers"
)

// APIHandler represents a lambda handler for API requests
type APIHandler struct {
	handleFunc interface{}
	logger     loggers.LoggerInterface
}

// NewAPIHandler creates a new APIHandler instance
func NewAPIHandler(handleFunc interface{}, logger loggers.LoggerInterface) *APIHandler {
	return &APIHandler{handleFunc, logger}
}

// Invoke handles a request
func (handler *APIHandler) Invoke(context context.Context, requestBytes []byte) (responseBytes []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			responseBytes, err = handler.internalError("Unexpected panic.", err)
		}
	}()
	if handler.handleFunc == nil {
		return handler.internalError("Missing handle function.", nil)
	}
	handleFuncType := reflect.TypeOf(handler.handleFunc)
	if handleFuncType.Kind() != reflect.Func {
		return handler.internalError("Unexpected handle function type.", nil)
	}
	if handleFuncType.NumIn() != 1 {
		return handler.internalError("Unexpected number of handle function arguments.", nil)
	}
	request := reflect.New(handleFuncType.In(0))
	err = json.Unmarshal(requestBytes, request.Interface())
	if err != nil {
		return handler.badRequestError("Failed to decode request JSON.")
	}
	args := []reflect.Value{request.Elem()}
	returnValues := reflect.ValueOf(handler.handleFunc).Call(args)
	numReturnValues := len(returnValues)
	if !(numReturnValues == 1 || numReturnValues == 2) {
		return handler.internalError("Unexpected number of handle function return values.", nil)
	}
	var errorValue = returnValues[numReturnValues-1].Interface()
	var response interface{}
	if errorValue != nil {
		err, ok := errorValue.(error)
		if !ok {
			return handler.internalError("Unexpected non-error return value.", nil)
		}
		httpError, ok := err.(errors.HTTPError)
		if !ok {
			return handler.internalError("Handle function returned an unexpected error.", err)
		}
		response = httpError
	} else if numReturnValues == 2 {
		response = returnValues[0].Interface()
	}
	if response != nil {
		responseBytes, err = json.Marshal(response)
		if err != nil {
			return handler.internalError("Failed to encode response JSON.", err)
		}
	}
	return responseBytes, nil
}

func (handler APIHandler) badRequestError(message string) ([]byte, error) {
	response := errors.NewBadRequestError(message)
	return json.Marshal(response)
}

func (handler APIHandler) internalError(message string, err error) ([]byte, error) {
	if handler.logger != nil {
		handler.logger.Error(message, err)
	}
	response := errors.NewInternalServerError()
	return json.Marshal(response)
}
