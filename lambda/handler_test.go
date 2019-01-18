package lambda_test

import (
	"errors"
	"os"
	"testing"

	coreErrors "github.com/bitterbattles/api/core/errors"
	"github.com/bitterbattles/api/core/mocks"
	"github.com/bitterbattles/api/core/tests"
	"github.com/bitterbattles/api/lambda"
)

var logger *mocks.Logger
var requestJSON = `{"key1":1,"key2":"value2","key3":true}`
var internalError string = `{"statusCode":500,"message":"Something unexpected happened. Please try again later."}`

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

func handleRequest(request request) (response, error) {
	return response{false, "val2", 3}, nil
}

func TestMain(m *testing.M) {
	logger = mocks.NewLogger()
	os.Exit(m.Run())
}

func TestNilFunction(t *testing.T) {
	reset()
	handler := lambda.NewHandler(nil, logger)
	response := invoke(t, handler, requestJSON)
	tests.AssertEquals(t, response, internalError)
	message, err := logger.GetLastEntry()
	tests.AssertNotNil(t, message)
	tests.AssertNil(t, err)
}

func TestBadFunctionType(t *testing.T) {
	reset()
	handler := lambda.NewHandler(1, logger)
	response := invoke(t, handler, requestJSON)
	tests.AssertEquals(t, response, internalError)
	message, err := logger.GetLastEntry()
	tests.AssertNotNil(t, message)
	tests.AssertNil(t, err)
}

func TestTooManyInputParams(t *testing.T) {
	reset()
	function := func(a int, b int) {}
	handler := lambda.NewHandler(function, logger)
	response := invoke(t, handler, requestJSON)
	tests.AssertEquals(t, response, internalError)
	message, err := logger.GetLastEntry()
	tests.AssertNotNil(t, message)
	tests.AssertNil(t, err)
}

func TestBadRequestJSON(t *testing.T) {
	reset()
	handler := lambda.NewHandler(handleRequest, logger)
	response := invoke(t, handler, "notjson")
	tests.AssertEquals(t, response, internalError)
	message, err := logger.GetLastEntry()
	tests.AssertNotNil(t, message)
	tests.AssertNotNil(t, err)
}

func TestTooFewReturnParams(t *testing.T) {
	reset()
	function := func(request request) {}
	handler := lambda.NewHandler(function, logger)
	response := invoke(t, handler, requestJSON)
	tests.AssertEquals(t, response, internalError)
	message, err := logger.GetLastEntry()
	tests.AssertNotNil(t, message)
	tests.AssertNil(t, err)
}

func TestBadErrorType(t *testing.T) {
	reset()
	function := func(request request) response { return response{} }
	handler := lambda.NewHandler(function, logger)
	response := invoke(t, handler, requestJSON)
	tests.AssertEquals(t, response, internalError)
	message, err := logger.GetLastEntry()
	tests.AssertNotNil(t, message)
	tests.AssertNil(t, err)
}

func TestUnepectedError(t *testing.T) {
	reset()
	expectedErr := errors.New("error")
	function := func(request request) (response, error) { return response{}, expectedErr }
	handler := lambda.NewHandler(function, logger)
	response := invoke(t, handler, requestJSON)
	tests.AssertEquals(t, response, internalError)
	message, err := logger.GetLastEntry()
	tests.AssertNotNil(t, message)
	tests.AssertEquals(t, err, expectedErr)
}

func TestExpectedError(t *testing.T) {
	reset()
	function := func(request request) (response, error) {
		return response{}, coreErrors.NewBadRequestError("badrequest")
	}
	handler := lambda.NewHandler(function, logger)
	response := invoke(t, handler, requestJSON)
	expected := `{"statusCode":400,"message":"badrequest"}`
	tests.AssertEquals(t, response, expected)
	message, err := logger.GetLastEntry()
	tests.AssertNil(t, message)
	tests.AssertNil(t, err)
}

func TestPanic(t *testing.T) {
	reset()
	expectedErr := errors.New("error")
	function := func(request request) (response, error) {
		panic(expectedErr)
	}
	handler := lambda.NewHandler(function, logger)
	response := invoke(t, handler, requestJSON)
	tests.AssertEquals(t, response, internalError)
	message, err := logger.GetLastEntry()
	tests.AssertNil(t, message)
	tests.AssertEquals(t, err, expectedErr)
}

func TestSuccess(t *testing.T) {
	reset()
	handler := lambda.NewHandler(handleRequest, logger)
	response := invoke(t, handler, requestJSON)
	expected := `{"key1":false,"key2":"val2","key3":3}`
	tests.AssertEquals(t, response, expected)
	message, err := logger.GetLastEntry()
	tests.AssertNotNil(t, message)
	tests.AssertNil(t, err)
}

func invoke(t *testing.T, handler *lambda.Handler, requestJSON string) string {
	requestBytes := []byte(requestJSON)
	responseBytes, err := handler.Invoke(nil, requestBytes)
	tests.AssertNil(t, err)
	response := string(responseBytes)
	return response
}

func reset() {
	logger.Reset()
}
