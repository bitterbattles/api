package handlers_test

import (
	"errors"
	"testing"

	. "github.com/bitterbattles/api/common/handlers"
	. "github.com/bitterbattles/api/common/tests"
)

const dynamoEventJSON = `{"records":[{"eventID":"id1"},{"eventID":"id2"}]}`

type testStreamHandler struct {
	returnError bool
	value       string
}

func (handler *testStreamHandler) Handle(event *DynamoEvent) error {
	if handler.returnError {
		return errors.New("error")
	}
	handler.value = event.Records[1].EventID
	return nil
}

func TestStreamHandlerBadJSON(t *testing.T) {
	subHandler := testStreamHandler{}
	handler := NewStreamHandler(&subHandler)
	err := invokeStreamHandler(t, handler, "notjson")
	AssertNotNil(t, err)
	AssertEquals(t, subHandler.value, "")
}

func TestStreamHandlerError(t *testing.T) {
	subHandler := testStreamHandler{returnError: true}
	handler := NewStreamHandler(&subHandler)
	err := invokeStreamHandler(t, handler, dynamoEventJSON)
	AssertNotNil(t, err)
	AssertEquals(t, err.Error(), "error")
}

func TestStreamHandlerSuccess(t *testing.T) {
	subHandler := testStreamHandler{}
	handler := NewStreamHandler(&subHandler)
	err := invokeStreamHandler(t, handler, dynamoEventJSON)
	AssertNil(t, err)
	AssertEquals(t, subHandler.value, "id2")
}

func invokeStreamHandler(t *testing.T, handler *StreamHandler, eventJSON string) error {
	eventBytes := []byte(eventJSON)
	resultBytes, err := handler.Invoke(nil, eventBytes)
	AssertNil(t, resultBytes)
	return err
}
