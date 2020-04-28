package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/pkg/aws"
	"github.com/bitterbattles/api/pkg/comments"
	"github.com/bitterbattles/api/pkg/dynamo"
	"github.com/bitterbattles/api/pkg/elasticsearch"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/users"
)

func main() {
	esClient := elasticsearch.NewClient(os.Getenv("ELASTICSEARCH_ADDRESS"))
	commentsIndex := comments.NewIndex(esClient)
	session := aws.NewSession(os.Getenv("AWS_REGION"))
	dynamoClient := dynamo.NewClient(session)
	commentsRepository := comments.NewRepository(dynamoClient)
	usersRepository := users.NewRepository(dynamoClient)
	processor := NewProcessor(commentsIndex, commentsRepository, usersRepository)
	handler := api.NewHandler(true, os.Getenv("ACCESS_TOKEN_SECRET"), processor)
	lambda.StartHandler(handler)
}
