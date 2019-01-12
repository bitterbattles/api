package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/battles"
)

func main() {
	index := battles.NewIndex()
	table := battles.NewTable()
	manager := battles.NewManager(index, table)
	controller := battles.NewController(manager)
	lambda.Start(controller.HandleGet)
}
