package handlers_test

import (
	"errors"
	"testing"

	commonErrors "github.com/bitterbattles/api/common/errors"
	. "github.com/bitterbattles/api/common/handlers"
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
	handler := NewAPIHandler(nil)
	testAPIHandler(t, handler, requestJSON, internalErrorResponse)
}

func TestAPIHandlerBadFunctionType(t *testing.T) {
	handler := NewAPIHandler(1)
	testAPIHandler(t, handler, requestJSON, internalErrorResponse)
}

func TestAPIHandlerTooManyInputParams(t *testing.T) {
	function := func(a int, b int) {}
	handler := NewAPIHandler(function)
	testAPIHandler(t, handler, requestJSON, internalErrorResponse)
}

func TestAPIHandlerBadRequestJSON(t *testing.T) {
	handler := NewAPIHandler(apiHandlerFunction)
	expectedResponse := `{"statusCode":400,"message":"Failed to decode request JSON."}`
	testAPIHandler(t, handler, "notjson", expectedResponse)
}

func TestAPIHandlerTooFewReturnParams(t *testing.T) {
	function := func(request request) {}
	handler := NewAPIHandler(function)
	testAPIHandler(t, handler, requestJSON, internalErrorResponse)
}

func TestAPIHandlerBadErrorType(t *testing.T) {
	function := func(request request) response { return response{} }
	handler := NewAPIHandler(function)
	testAPIHandler(t, handler, requestJSON, internalErrorResponse)
}

func TestAPIHandlerUnexpectedError(t *testing.T) {
	expectedErr := errors.New("error")
	function := func(request request) (response, error) { return response{}, expectedErr }
	handler := NewAPIHandler(function)
	testAPIHandler(t, handler, requestJSON, internalErrorResponse)
}

func TestAPIHandlerExpectedError(t *testing.T) {
	function := func(request request) (response, error) {
		return response{}, commonErrors.NewBadRequestError("badrequest")
	}
	handler := NewAPIHandler(function)
	expectedResponse := `{"statusCode":400,"message":"badrequest"}`
	testAPIHandler(t, handler, requestJSON, expectedResponse)
}

func TestAPIHandlerPanic(t *testing.T) {
	expectedErr := errors.New("error")
	function := func(request request) (response, error) {
		panic(expectedErr)
	}
	handler := NewAPIHandler(function)
	testAPIHandler(t, handler, requestJSON, internalErrorResponse)
}

func TestAPIHandlerSuccess(t *testing.T) {
	handler := NewAPIHandler(apiHandlerFunction)
	expectedResponse := `{"key1":false,"key2":"val2","key3":3}`
	testAPIHandler(t, handler, requestJSON, expectedResponse)
}

func testAPIHandler(t *testing.T, handler *APIHandler, requestJSON string, expectedResponse string) {
	requestBytes := []byte(requestJSON)
	responseBytes, err := handler.Invoke(nil, requestBytes)
	AssertNil(t, err)
	response := string(responseBytes)
	AssertEquals(t, response, expectedResponse)
}
