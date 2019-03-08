package stream_test

import (
	"errors"
	"testing"

	. "github.com/bitterbattles/api/pkg/lambda/stream"
	. "github.com/bitterbattles/api/pkg/tests"
)

const eventJSON = `{"records":[{"eventID":"id1"},{"eventID":"id2"}]}`

type testProcessor struct {
	returnError bool
	value       string
}

func (processor *testProcessor) Process(event *Event) error {
	if processor.returnError {
		return errors.New("error")
	}
	processor.value = event.Records[1].EventID
	return nil
}

func TestHandlerBadJSON(t *testing.T) {
	processor := testProcessor{}
	handler := NewHandler(&processor)
	err := testHandler(t, handler, "notjson")
	AssertNotNil(t, err)
	AssertEquals(t, processor.value, "")
}

func TestHandlerError(t *testing.T) {
	processor := testProcessor{returnError: true}
	handler := NewHandler(&processor)
	err := testHandler(t, handler, eventJSON)
	AssertNotNil(t, err)
	AssertEquals(t, err.Error(), "error")
}

func TestHandlerSuccess(t *testing.T) {
	processor := testProcessor{}
	handler := NewHandler(&processor)
	err := testHandler(t, handler, eventJSON)
	AssertNil(t, err)
	AssertEquals(t, processor.value, "id2")
}

func testHandler(t *testing.T, handler *Handler, eventJSON string) error {
	eventBytes := []byte(eventJSON)
	resultBytes, err := handler.Invoke(nil, eventBytes)
	AssertNil(t, resultBytes)
	return err
}
