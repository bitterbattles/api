package api_test

import (
	"encoding/json"
	"errors"
	"testing"

	commonErrors "github.com/bitterbattles/api/pkg/common/errors"
	"github.com/bitterbattles/api/pkg/common/http"
	. "github.com/bitterbattles/api/pkg/common/lambda/api"
	. "github.com/bitterbattles/api/pkg/common/tests"
)

var headers = map[string]string{
	"Content-Type": "application/json",
}

type testHandler struct {
	returnHTTPError       bool
	returnUnexpectedError bool
	returnNilResponse     bool
	panic                 bool
}

func (handler *testHandler) Handle(request *http.Request) (*http.Response, error) {
	if handler.returnHTTPError {
		return nil, commonErrors.NewBadRequestError("Bad request.")
	}
	if handler.returnUnexpectedError {
		return nil, errors.New("unexpected")
	}
	if handler.returnNilResponse {
		return nil, nil
	}
	if handler.panic {
		panic(errors.New("panic"))
	}
	response := &http.Response{
		StatusCode: http.OK,
		Body:       `{"key1":"value1","key2":2,"key3":true}`,
	}
	return response, nil
}

func TestHandlerBadRequestJSON(t *testing.T) {
	handler := NewHandler(&testHandler{})
	responseBytes, err := handler.Invoke(nil, []byte("notjson"))
	AssertNotNil(t, responseBytes)
	AssertNil(t, err)
	response := http.Response{}
	err = json.Unmarshal(responseBytes, &response)
	AssertNil(t, err)
	AssertEquals(t, response.StatusCode, http.InternalServerError)
}

func TestHandlerBadContentType(t *testing.T) {
	handler := NewHandler(&testHandler{})
	request := &http.Request{Body: "body"}
	expectedBody := `{"errorMessage":"Only JSON content type is accepted."}`
	invokeHandler(t, handler, request, http.UnsupportedMediaType, expectedBody)
}

func TestHandlerHandledError(t *testing.T) {
	handler := NewHandler(&testHandler{returnHTTPError: true})
	request := &http.Request{Headers: headers}
	expectedBody := `{"errorMessage":"Bad request."}`
	invokeHandler(t, handler, request, http.BadRequest, expectedBody)
}

func TestHandlerUnhandledError(t *testing.T) {
	handler := NewHandler(&testHandler{returnUnexpectedError: true})
	request := &http.Request{Headers: headers}
	expectedBody := `{"errorMessage":"Something unexpected happened. Please try again later."}`
	invokeHandler(t, handler, request, http.InternalServerError, expectedBody)
}

func TestHandlerMissingResposne(t *testing.T) {
	handler := NewHandler(&testHandler{returnNilResponse: true})
	request := &http.Request{Headers: headers}
	expectedBody := `{"errorMessage":"Something unexpected happened. Please try again later."}`
	invokeHandler(t, handler, request, http.InternalServerError, expectedBody)
}

func TestHandlerPanic(t *testing.T) {
	handler := NewHandler(&testHandler{panic: true})
	request := &http.Request{Headers: headers}
	expectedBody := `{"errorMessage":"Something unexpected happened. Please try again later."}`
	invokeHandler(t, handler, request, http.InternalServerError, expectedBody)
}

func TestHandlerSuccess(t *testing.T) {
	handler := NewHandler(&testHandler{})
	request := &http.Request{Headers: headers}
	expectedBody := `{"key1":"value1","key2":2,"key3":true}`
	invokeHandler(t, handler, request, http.OK, expectedBody)
}

func invokeHandler(t *testing.T, handler *Handler, request *http.Request, expectedStatusCode int, expectedBody string) {
	requestBytes, _ := json.Marshal(request)
	responseBytes, err := handler.Invoke(nil, requestBytes)
	AssertNotNil(t, responseBytes)
	AssertNil(t, err)
	response := http.Response{}
	err = json.Unmarshal(responseBytes, &response)
	AssertNil(t, err)
	AssertEquals(t, response.StatusCode, expectedStatusCode)
	AssertEquals(t, response.Headers[http.ContentTypeHeader], http.ApplicationJSON)
	AssertEquals(t, response.Body, expectedBody)
}
