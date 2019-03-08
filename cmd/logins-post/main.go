package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/pkg/aws"
	"github.com/bitterbattles/api/pkg/dynamo"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/users"
)

func main() {
	session := aws.NewSession(os.Getenv("AWS_REGION"))
	dynamoClient := dynamo.NewClient(session)
	repository := users.NewRepository(dynamoClient)
	tokenSecret := os.Getenv("TOKEN_SECRET")
	processor := NewProcessor(repository, tokenSecret)
	handler := api.NewHandler(false, tokenSecret, processor)
	lambda.StartHandler(handler)
}
