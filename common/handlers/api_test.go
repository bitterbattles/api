package handlers_test

import (
	"errors"
	"testing"

	commonErrors "github.com/bitterbattles/api/common/errors"
	. "github.com/bitterbattles/api/common/handlers"
	"github.com/bitterbattles/api/common/loggers/mocks"
	. "github.com/bitterbattles/api/common/tests"
)

const requestJSON = `{"key1":1,"key2":"value2","key3":true}`
const internalErrorResponse = `{"statusCode":500,"message":"Something unexpected happened. Please try again later."}`

type request struct {
	key1 int
	key2 string
	key3 bool
}

type response struct {
	Key1 bool   `json:"key1"`
	Key2 string `json:"key2"`
	Key3 int    `json:"key3"`
}

func apiHandlerFunction(request *request) (*response, error) {
	return &response{false, "val2", 3}, nil
}

func TestAPIHandlerNilFunction(t *testing.T) {
	logger := mocks.NewLogger()
	handler := NewAPIHandler(nil, logger)
	testAPIHandler(t, handler, requestJSON, internalErrorResponse)
	validateLogger(t, logger, true, false, nil)
}

func TestAPIHandlerBadFunctionType(t *testing.T) {
	logger := mocks.NewLogger()
	handler := NewAPIHandler(1, logger)
	testAPIHandler(t, handler, requestJSON, internalErrorResponse)
	validateLogger(t, logger, true, false, nil)
}

func TestAPIHandlerTooManyInputParams(t *testing.T) {
	logger := mocks.NewLogger()
	function := func(a int, b int) {}
	handler := NewAPIHandler(function, logger)
	testAPIHandler(t, handler, requestJSON, internalErrorResponse)
	validateLogger(t, logger, true, false, nil)
}

func TestAPIHandlerBadRequestJSON(t *testing.T) {
	logger := mocks.NewLogger()
	handler := NewAPIHandler(apiHandlerFunction, logger)
	expectedResponse := `{"statusCode":400,"message":"Failed to decode request JSON."}`
	testAPIHandler(t, handler, "notjson", expectedResponse)
	validateLogger(t, logger, true, false, nil)
}

func TestAPIHandlerTooFewReturnParams(t *testing.T) {
	logger := mocks.NewLogger()
	function := func(request request) {}
	handler := NewAPIHandler(function, logger)
	testAPIHandler(t, handler, requestJSON, internalErrorResponse)
	validateLogger(t, logger, true, false, nil)
}

func TestAPIHandlerBadErrorType(t *testing.T) {
	logger := mocks.NewLogger()
	function := func(request request) response { return response{} }
	handler := NewAPIHandler(function, logger)
	testAPIHandler(t, handler, requestJSON, internalErrorResponse)
	validateLogger(t, logger, true, false, nil)
}

func TestAPIHandlerUnexpectedError(t *testing.T) {
	logger := mocks.NewLogger()
	expectedErr := errors.New("error")
	function := func(request request) (response, error) { return response{}, expectedErr }
	handler := NewAPIHandler(function, logger)
	testAPIHandler(t, handler, requestJSON, internalErrorResponse)
	validateLogger(t, logger, true, true, expectedErr)
}

func TestAPIHandlerExpectedError(t *testing.T) {
	logger := mocks.NewLogger()
	function := func(request request) (response, error) {
		return response{}, commonErrors.NewBadRequestError("badrequest")
	}
	handler := NewAPIHandler(function, logger)
	expectedResponse := `{"statusCode":400,"message":"badrequest"}`
	testAPIHandler(t, handler, requestJSON, expectedResponse)
	validateLogger(t, logger, false, false, nil)
}

func TestAPIHandlerPanic(t *testing.T) {
	logger := mocks.NewLogger()
	expectedErr := errors.New("error")
	function := func(request request) (response, error) {
		panic(expectedErr)
	}
	handler := NewAPIHandler(function, logger)
	testAPIHandler(t, handler, requestJSON, internalErrorResponse)
	validateLogger(t, logger, false, true, expectedErr)
}

func TestAPIHandlerSuccess(t *testing.T) {
	logger := mocks.NewLogger()
	handler := NewAPIHandler(apiHandlerFunction, logger)
	expectedResponse := `{"key1":false,"key2":"val2","key3":3}`
	testAPIHandler(t, handler, requestJSON, expectedResponse)
	validateLogger(t, logger, false, false, nil)
}

func testAPIHandler(t *testing.T, handler *APIHandler, requestJSON string, expectedResponse string) {
	requestBytes := []byte(requestJSON)
	responseBytes, err := handler.Invoke(nil, requestBytes)
	AssertNil(t, err)
	response := string(responseBytes)
	AssertEquals(t, response, expectedResponse)
}

func validateLogger(t *testing.T, logger *mocks.Logger, nonNilMessage bool, nonNilError bool, expectedError error) {
	message, err := logger.GetLastEntry()
	if nonNilMessage {
		AssertNotNil(t, message)
	} else {
		AssertNil(t, message)
	}
	if nonNilError {
		AssertNotNil(t, err)
		AssertEquals(t, err, expectedError)
	} else {
		AssertNil(t, err)
	}
}
