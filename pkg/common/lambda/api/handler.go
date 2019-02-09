package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/bitterbattles/api/pkg/common/errors"
	"github.com/bitterbattles/api/pkg/common/http"
)

// HandlerInterface defines an interface for handling API Gateway proxy requests
type HandlerInterface interface {
	Handle(*http.Request) (*http.Response, error)
}

// Handler represents a lambda handler for API Gateway proxy requests
type Handler struct {
	HandlerInterface
}

// NewHandler creates a new Handler instance
func NewHandler(handler HandlerInterface) *Handler {
	return &Handler{handler}
}

// Invoke handles a request
func (handler *Handler) Invoke(context context.Context, requestBytes []byte) (responseBytes []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			responseBytes, err = handler.unhandledError("Unexpected panic.", err)
		}
	}()
	request := http.Request{}
	err = json.Unmarshal(requestBytes, &request)
	if err != nil {
		return handler.unhandledError("Failed to decode proxy request JSON.", err)
	}
	contentType := request.Headers[http.ContentTypeHeader]
	if contentType != http.ApplicationJSON {
		return handler.handledError(http.UnsupportedMediaType, "Only JSON content type is accepted.")
	}
	response, err := handler.Handle(&request)
	if err != nil {
		httpError, ok := err.(errors.HTTPError)
		if ok {
			return handler.handledError(httpError.StatusCode, httpError.Message)
		}
		return handler.unhandledError("Handle function returned an unexpected error.", err)
	}
	if response == nil {
		return handler.unhandledError("Handle function returned nil response.", nil)
	}
	if response.Headers == nil {
		response.Headers = make(map[string]string)
	}
	response.Headers[http.ContentTypeHeader] = http.ApplicationJSON
	responseBytes, err = json.Marshal(response)
	if err != nil {
		return handler.unhandledError("Failed to encode proxy response JSON.", err)
	}
	return responseBytes, nil
}

func (handler *Handler) handledError(statusCode int, message string) ([]byte, error) {
	headers := map[string]string{
		http.ContentTypeHeader: http.ApplicationJSON,
	}
	response := &http.Response{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       fmt.Sprintf(`{"errorMessage":"%s"}`, message),
	}
	return json.Marshal(response)
}

func (handler *Handler) unhandledError(message string, err error) ([]byte, error) {
	log.Println(message, "Error:", err)
	return handler.handledError(http.InternalServerError, "Something unexpected happened. Please try again later.")
}
