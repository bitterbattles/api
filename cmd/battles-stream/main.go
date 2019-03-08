package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/pkg/lambda/stream"
	"github.com/bitterbattles/api/pkg/ranks"
	"github.com/bitterbattles/api/pkg/redis"
)

func main() {
	redisClient := redis.NewClient(os.Getenv("REDIS_ADDRESS"))
	repository := ranks.NewRepository(redisClient)
	processor := NewProcessor(repository)
	handler := stream.NewHandler(processor)
	lambda.StartHandler(handler)
}
