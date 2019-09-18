package main

import (
	"os"

	"github.com/bitterbattles/api/pkg/users"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/pkg/aws"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/dynamo"
	"github.com/bitterbattles/api/pkg/index"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/redis"
)

func main() {
	redisClient := redis.NewClient(os.Getenv("REDIS_ADDRESS"))
	indexRepository := index.NewRepository(redisClient)
	indexer := battles.NewIndexer(indexRepository)
	session := aws.NewSession(os.Getenv("AWS_REGION"))
	dynamoClient := dynamo.NewClient(session)
	battlesRepository := battles.NewRepository(dynamoClient)
	usersRepository := users.NewRepository(dynamoClient)
	processor := NewProcessor(indexer, battlesRepository, usersRepository)
	handler := api.NewHandler(true, os.Getenv("ACCESS_TOKEN_SECRET"), processor)
	lambda.StartHandler(handler)
}
