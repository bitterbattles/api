package stream_test

import (
	"errors"
	"testing"

	. "github.com/bitterbattles/api/pkg/common/lambda/stream"
	. "github.com/bitterbattles/api/pkg/common/tests"
)

const eventJSON = `{"records":[{"eventID":"id1"},{"eventID":"id2"}]}`

type testHandler struct {
	returnError bool
	value       string
}

func (handler *testHandler) Handle(event *Event) error {
	if handler.returnError {
		return errors.New("error")
	}
	handler.value = event.Records[1].EventID
	return nil
}

func TestHandlerBadJSON(t *testing.T) {
	subHandler := testHandler{}
	handler := NewHandler(&subHandler)
	err := invokeHandler(t, handler, "notjson")
	AssertNotNil(t, err)
	AssertEquals(t, subHandler.value, "")
}

func TestHandlerError(t *testing.T) {
	subHandler := testHandler{returnError: true}
	handler := NewHandler(&subHandler)
	err := invokeHandler(t, handler, eventJSON)
	AssertNotNil(t, err)
	AssertEquals(t, err.Error(), "error")
}

func TestHandlerSuccess(t *testing.T) {
	subHandler := testHandler{}
	handler := NewHandler(&subHandler)
	err := invokeHandler(t, handler, eventJSON)
	AssertNil(t, err)
	AssertEquals(t, subHandler.value, "id2")
}

func invokeHandler(t *testing.T, handler *Handler, eventJSON string) error {
	eventBytes := []byte(eventJSON)
	resultBytes, err := handler.Invoke(nil, eventBytes)
	AssertNil(t, resultBytes)
	return err
}
