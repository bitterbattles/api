package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/pkg/aws"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/dynamo"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/votes"
)

func main() {
	session := aws.NewSession(os.Getenv("AWS_REGION"))
	dynamoClient := dynamo.NewClient(session)
	battlesRepository := battles.NewRepository(dynamoClient)
	votesRepository := votes.NewRepository(dynamoClient)
	processor := NewProcessor(battlesRepository, votesRepository)
	handler := api.NewHandler(false, os.Getenv("ACCESS_TOKEN_SECRET"), processor)
	lambda.StartHandler(handler)
}
