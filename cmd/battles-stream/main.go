package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/index"
	"github.com/bitterbattles/api/pkg/lambda/stream"
	"github.com/bitterbattles/api/pkg/redis"
)

func main() {
	redisClient := redis.NewClient(os.Getenv("REDIS_ADDRESS"))
	repository := index.NewRepository(redisClient)
	indexer := battles.NewIndexer(repository)
	processor := NewProcessor(indexer)
	handler := stream.NewHandler(processor)
	lambda.StartHandler(handler)
}
