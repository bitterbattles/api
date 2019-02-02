package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/pkg/common/bootstrap"
	"github.com/bitterbattles/api/pkg/votes"
)

func main() {
	session := bootstrap.NewSession()
	dynamoClient := bootstrap.NewDynamoClient(session)
	repository := votes.NewRepository(dynamoClient)
	handler := NewHandler(repository)
	lambda.StartHandler(handler)
}
