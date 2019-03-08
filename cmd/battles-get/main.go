package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/pkg/aws"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/dynamo"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/ranks"
	"github.com/bitterbattles/api/pkg/redis"
)

func main() {
	session := aws.NewSession(os.Getenv("AWS_REGION"))
	dynamoClient := dynamo.NewClient(session)
	battlesRepository := battles.NewRepository(dynamoClient)
	redisClient := redis.NewClient(os.Getenv("REDIS_ADDRESS"))
	ranksRepository := ranks.NewRepository(redisClient)
	processor := NewProcessor(battlesRepository, ranksRepository)
	handler := api.NewHandler(false, os.Getenv("TOKEN_SECRET"), processor)
	lambda.StartHandler(handler)
}
