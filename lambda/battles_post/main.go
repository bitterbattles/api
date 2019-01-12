package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/battles"
	"github.com/bitterbattles/api/mocks"
)

func main() {
	index := mocks.NewMockIndex()
	table := mocks.NewMockTable()
	manager := battles.NewManager(index, table)
	controller := battles.NewController(manager)
	lambda.Start(controller.HandlePost)
}
