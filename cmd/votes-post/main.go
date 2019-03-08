package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/pkg/aws"
	"github.com/bitterbattles/api/pkg/dynamo"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/votes"
)

func main() {
	session := aws.NewSession(os.Getenv("AWS_REGION"))
	dynamoClient := dynamo.NewClient(session)
	repository := votes.NewRepository(dynamoClient)
	processor := NewProcessor(repository)
	handler := api.NewHandler(true, os.Getenv("TOKEN_SECRET"), processor)
	lambda.StartHandler(handler)
}
