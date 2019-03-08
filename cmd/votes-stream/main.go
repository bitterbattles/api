package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/pkg/aws"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/dynamo"
	"github.com/bitterbattles/api/pkg/lambda/stream"
)

func main() {
	session := aws.NewSession(os.Getenv("AWS_REGION"))
	dynamoClient := dynamo.NewClient(session)
	repository := battles.NewRepository(dynamoClient)
	processor := NewProcessor(repository)
	handler := stream.NewHandler(processor)
	lambda.StartHandler(handler)
}
