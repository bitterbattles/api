package api_test

import (
	"encoding/json"
	"errors"
	"testing"

	apiErrors "github.com/bitterbattles/api/pkg/errors"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
)

var testHeaders = map[string]string{
	"authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VySWQiLCJpYXQiOjE1NTE5MjA0OTIsImV4cCI6NDc1NTU1MjgwN30.RDlto1SRXTsX1Gvi3vTaJs5itLadm2KAu1vQVFoqPuw",
	"content-type":  "application/json",
}

type testRequest struct {
	Key1 string `json:"key1"`
	Key2 int    `json:"key2"`
	Key3 bool   `json:"key3"`
}

type testResponse struct {
	Key4 string `json:"key4"`
	Key5 int    `json:"key5"`
	Key6 bool   `json:"key6"`
}

type testProcessor struct {
	requireBody           bool
	returnHTTPError       bool
	returnUnexpectedError bool
	returnNilResponse     bool
	panic                 bool
}

func (processor *testProcessor) NewRequestBody() interface{} {
	if processor.requireBody {
		return &testRequest{}
	}
	return nil
}

func (processor *testProcessor) Process(input *api.Input) (*api.Output, error) {
	if processor.returnHTTPError {
		return nil, apiErrors.NewBadRequestError("Bad request.")
	}
	if processor.returnUnexpectedError {
		return nil, errors.New("unexpected")
	}
	if processor.returnNilResponse {
		return nil, nil
	}
	if processor.panic {
		panic(errors.New("panic"))
	}
	var response *testResponse
	if processor.requireBody {
		request, _ := input.RequestBody.(*testRequest)
		response = &testResponse{
			Key4: request.Key1,
			Key5: request.Key2,
			Key6: request.Key3,
		}
	} else {
		response = &testResponse{
			Key4: "value4",
			Key5: 5,
			Key6: true,
		}
	}
	output := api.NewOutput(response)
	return output, nil
}

func TestHandlerBadRequestJSON(t *testing.T) {
	processor := &testProcessor{}
	handler := NewHandler(false, "", processor)
	responseBytes, err := handler.Invoke(nil, []byte("notjson"))
	AssertNotNil(t, responseBytes)
	AssertNil(t, err)
	response := http.Response{}
	err = json.Unmarshal(responseBytes, &response)
	AssertNil(t, err)
	AssertEquals(t, response.StatusCode, http.InternalServerError)
}

func TestHandlerMissingAuth(t *testing.T) {
	processor := &testProcessor{}
	handler := NewHandler(true, "", processor)
	request := &http.Request{}
	expectedBody := `{"errorCode":401,"errorMessage":"Authorization is required."}`
	testHandler(t, handler, request, http.Unauthorized, expectedBody)
}

func TestHandlerMalformedAuth(t *testing.T) {
	processor := &testProcessor{}
	handler := NewHandler(true, "", processor)
	headers := map[string]string{
		"Authorization": "invalid",
	}
	request := &http.Request{Headers: headers}
	expectedBody := `{"errorCode":403,"errorMessage":"Malformed Authorization header."}`
	testHandler(t, handler, request, http.Forbidden, expectedBody)
}

func TestHandlerBadToken(t *testing.T) {
	processor := &testProcessor{}
	handler := NewHandler(true, "tokenSecret", processor)
	headers := map[string]string{
		"Authorization": "Bearer invalid",
	}
	request := &http.Request{Headers: headers}
	expectedBody := `{"errorCode":403,"errorMessage":"Invalid token."}`
	testHandler(t, handler, request, http.Forbidden, expectedBody)
}

func TestHandlerExpiredToken(t *testing.T) {
	processor := &testProcessor{}
	handler := NewHandler(true, "tokenSecret", processor)
	headers := map[string]string{
		"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VySWQiLCJpYXQiOjE1NTE5MjA0OTIsImV4cCI6MTU1MTkyMDQ5M30.tLUDq1svsmPySqdg_pTjn68fD8ETjSDX626xi-CV5dM",
	}
	request := &http.Request{Headers: headers}
	expectedBody := `{"errorCode":403,"errorMessage":"Expired token."}`
	testHandler(t, handler, request, http.Forbidden, expectedBody)
}

func TestHandlerMissingBody(t *testing.T) {
	processor := &testProcessor{requireBody: true}
	handler := NewHandler(false, "tokenSecret", processor)
	request := &http.Request{}
	expectedBody := `{"errorCode":400,"errorMessage":"Request body is required."}`
	testHandler(t, handler, request, http.BadRequest, expectedBody)
}

func TestHandlerBadContentType(t *testing.T) {
	processor := &testProcessor{requireBody: true}
	handler := NewHandler(false, "tokenSecret", processor)
	request := &http.Request{Body: "body"}
	expectedBody := `{"errorCode":415,"errorMessage":"Only JSON content type is accepted."}`
	testHandler(t, handler, request, http.UnsupportedMediaType, expectedBody)
}

func TestHandlerBadBodyJSON(t *testing.T) {
	processor := &testProcessor{requireBody: true}
	handler := NewHandler(false, "tokenSecret", processor)
	request := &http.Request{Headers: testHeaders, Body: "body"}
	expectedBody := `{"errorCode":400,"errorMessage":"Failed to decode request body."}`
	testHandler(t, handler, request, http.BadRequest, expectedBody)
}

func TestHandlerExpectedError(t *testing.T) {
	processor := &testProcessor{returnHTTPError: true}
	handler := NewHandler(false, "tokenSecret", processor)
	request := &http.Request{Headers: testHeaders}
	expectedBody := `{"errorCode":400,"errorMessage":"Bad request."}`
	testHandler(t, handler, request, http.BadRequest, expectedBody)
}

func TestHandlerUnexpectedError(t *testing.T) {
	processor := &testProcessor{returnUnexpectedError: true}
	handler := NewHandler(false, "tokenSecret", processor)
	request := &http.Request{Headers: testHeaders}
	expectedBody := `{"errorCode":500,"errorMessage":"Something unexpected happened. Please try again later."}`
	testHandler(t, handler, request, http.InternalServerError, expectedBody)
}

func TestHandlerPanic(t *testing.T) {
	processor := &testProcessor{panic: true}
	handler := NewHandler(false, "tokenSecret", processor)
	request := &http.Request{Headers: testHeaders}
	expectedBody := `{"errorCode":500,"errorMessage":"Something unexpected happened. Please try again later."}`
	testHandler(t, handler, request, http.InternalServerError, expectedBody)
}

func TestHandlerMissingOutput(t *testing.T) {
	processor := &testProcessor{returnNilResponse: true}
	handler := NewHandler(false, "tokenSecret", processor)
	request := &http.Request{Headers: testHeaders}
	expectedBody := `{"errorCode":500,"errorMessage":"Something unexpected happened. Please try again later."}`
	testHandler(t, handler, request, http.InternalServerError, expectedBody)
}

func TestHandlerSuccess(t *testing.T) {
	processor := &testProcessor{requireBody: true}
	handler := NewHandler(true, "tokenSecret", processor)
	request := &http.Request{Headers: testHeaders, Body: `{"key1":"value1","key2":2,"key3":true}`}
	expectedBody := `{"key4":"value1","key5":2,"key6":true}`
	testHandler(t, handler, request, http.OK, expectedBody)
}

func testHandler(t *testing.T, handler *Handler, request *http.Request, expectedStatusCode int, expectedBody string) {
	requestBytes, _ := json.Marshal(request)
	responseBytes, err := handler.Invoke(nil, requestBytes)
	AssertNotNil(t, responseBytes)
	AssertNil(t, err)
	response := http.Response{}
	err = json.Unmarshal(responseBytes, &response)
	AssertNil(t, err)
	AssertEquals(t, response.StatusCode, expectedStatusCode)
	AssertEquals(t, response.Headers[http.AccessControlAllowOrigin], "*")
	AssertEquals(t, response.Headers[http.ContentType], http.ApplicationJSON)
	AssertEquals(t, response.Body, expectedBody)
}
