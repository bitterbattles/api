package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/common/bootstrap"
	"github.com/bitterbattles/api/pkg/ranks"
)

func main() {
	session := bootstrap.NewSession()
	dynamoClient := bootstrap.NewDynamoClient(session)
	redisClient := bootstrap.NewRedisClient()
	battlesRepository := battles.NewRepository(dynamoClient)
	ranksRepository := ranks.NewRepository(redisClient)
	handler := NewHandler(ranksRepository, battlesRepository)
	lambda.StartHandler(handler)
}
