package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/common/bootstrap"
)

func main() {
	session := bootstrap.NewSession()
	dynamoClient := bootstrap.NewDynamoClient(session)
	repository := battles.NewRepository(dynamoClient)
	handler := NewHandler(repository)
	lambda.StartHandler(handler)
}
