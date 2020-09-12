package api

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/bitterbattles/api/pkg/errors"
	"github.com/bitterbattles/api/pkg/http"
	"github.com/bitterbattles/api/pkg/jwt"
	"github.com/bitterbattles/api/pkg/time"
)

type errorResponse struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

// Handler represents a lambda handler for API Gateway proxy requests
type Handler struct {
	requiresAuth bool
	tokenSecret  string
	processor    ProcessorInterface
}

// NewHandler creates a new Handler instance
func NewHandler(requiresAuth bool, tokenSecret string, processor ProcessorInterface) *Handler {
	return &Handler{
		requiresAuth: requiresAuth,
		tokenSecret:  tokenSecret,
		processor:    processor,
	}
}

// Invoke handles a request
func (handler *Handler) Invoke(context context.Context, requestBytes []byte) (responseBytes []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			responseBytes, err = handler.newUnexpectedErrorResponse("Unexpected panic.", err)
		}
	}()
	request := &http.Request{}
	err = json.Unmarshal(requestBytes, request)
	if err != nil {
		return handler.newUnexpectedErrorResponse("Failed to decode proxy request.", err)
	}
	var authContext *AuthContext
	authorization := request.GetHeader(http.Authorization)
	if authorization != "" || handler.requiresAuth {
		if authorization == "" {
			return handler.newErrorResponse(http.Unauthorized, "Authorization is required.")
		}
		authParts := strings.Split(authorization, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			return handler.newErrorResponse(http.Forbidden, "Malformed Authorization header.")
		}
		authContext = &AuthContext{}
		err = jwt.DecodeHS256(authParts[1], handler.tokenSecret, authContext)
		if err != nil {
			return handler.newErrorResponse(http.Forbidden, "Invalid token.")
		}
		if authContext.ExpiresOn <= time.NowUnix() {
			return handler.newErrorResponse(http.Forbidden, "Expired token.")
		}
	}
	requestBody := handler.processor.NewRequestBody()
	if requestBody != nil {
		if request.Body == "" {
			return handler.newErrorResponse(http.BadRequest, "Request body is required.")
		}
		contentType := request.GetHeader(http.ContentType)
		if contentType != http.ApplicationJSON {
			return handler.newErrorResponse(http.UnsupportedMediaType, "Only JSON content type is accepted.")
		}
		err := json.Unmarshal([]byte(request.Body), requestBody)
		if err != nil {
			return handler.newErrorResponse(http.BadRequest, "Failed to decode request body.")
		}
	}
	input := &Input{
		PathParams:  request.PathParams,
		QueryParams: request.QueryParams,
		AuthContext: authContext,
		RequestBody: requestBody,
	}
	output, err := handler.processor.Process(input)
	if err != nil {
		httpError, ok := err.(*errors.HTTPError)
		if ok {
			return handler.newErrorResponseWithCode(httpError.StatusCode(), httpError.ErrorCode(), httpError.Error())
		}
		return handler.newUnexpectedErrorResponse("Processor returned an unexpected error.", err)
	}
	if output == nil {
		return handler.newUnexpectedErrorResponse("Processor returned nil output.", nil)
	}
	return handler.newResponse(output)
}

func (handler *Handler) newUnexpectedErrorResponse(message string, err error) ([]byte, error) {
	log.Println("Unexpected Error:", message, "-", err)
	return handler.newErrorResponse(http.InternalServerError, "Something unexpected happened. Please try again later.")
}

func (handler *Handler) newErrorResponse(statusCode int, errorMessage string) ([]byte, error) {
	return handler.newErrorResponseWithCode(statusCode, statusCode, errorMessage)
}

func (handler *Handler) newErrorResponseWithCode(statusCode int, errorCode int, errorMessage string) ([]byte, error) {
	output := &Output{
		StatusCode: statusCode,
		ResponseBody: &errorResponse{
			ErrorCode:    errorCode,
			ErrorMessage: errorMessage,
		},
	}
	return handler.newResponse(output)
}

func (handler *Handler) newResponse(output *Output) ([]byte, error) {
	var headers map[string]string
	var body string
	if output.ResponseBody != nil {
		bodyBytes, err := json.Marshal(output.ResponseBody)
		if err != nil {
			return nil, err
		}
		body = string(bodyBytes)
		headers = make(map[string]string)
		headers[http.AccessControlAllowOrigin] = "*"
		headers[http.ContentType] = http.ApplicationJSON
	}
	response := &http.Response{
		StatusCode: output.StatusCode,
		Headers:    headers,
		Body:       body,
	}
	return json.Marshal(response)
}
