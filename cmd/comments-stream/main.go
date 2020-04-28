package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/pkg/aws"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/comments"
	"github.com/bitterbattles/api/pkg/dynamo"
	"github.com/bitterbattles/api/pkg/elasticsearch"
	"github.com/bitterbattles/api/pkg/lambda/stream"
)

func main() {
	esClient := elasticsearch.NewClient(os.Getenv("ELASTICSEARCH_ADDRESS"))
	commentsIndex := comments.NewIndex(esClient)
	session := aws.NewSession(os.Getenv("AWS_REGION"))
	dynamoClient := dynamo.NewClient(session)
	battlesRepository := battles.NewRepository(dynamoClient)
	processor := NewProcessor(commentsIndex, battlesRepository)
	handler := stream.NewHandler(processor)
	lambda.StartHandler(handler)
}
