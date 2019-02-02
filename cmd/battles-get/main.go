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
	battlesRepository := battles.NewRepository(dynamoClient)
	ranksRepository := ranks.NewRepository(dynamoClient)
	handler := NewHandler(ranksRepository, battlesRepository)
	lambda.StartHandler(handler)
}
