package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/pkg/aws"
	"github.com/bitterbattles/api/pkg/comments"
	"github.com/bitterbattles/api/pkg/dynamo"
	"github.com/bitterbattles/api/pkg/index"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/redis"
)

func main() {
	redisClient := redis.NewClient(os.Getenv("REDIS_ADDRESS"))
	indexRepository := index.NewRepository(redisClient)
	indexer := comments.NewIndexer(indexRepository)
	session := aws.NewSession(os.Getenv("AWS_REGION"))
	dynamoClient := dynamo.NewClient(session)
	commentsRepository := comments.NewRepository(dynamoClient)
	processor := NewProcessor(indexer, commentsRepository)
	handler := api.NewHandler(false, os.Getenv("ACCESS_TOKEN_SECRET"), processor)
	lambda.StartHandler(handler)
}
