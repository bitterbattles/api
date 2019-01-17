package main

import (
	"github.com/bitterbattles/api/battles"
	"github.com/bitterbattles/api/lambda"
)

func main() {
	index := battles.NewIndex()
	repository := battles.NewRepository()
	manager := battles.NewManager(index, repository)
	controller := battles.NewController(manager)
	handler := lambda.NewHandler(controller.HandleGet)
	handler.Run()
}
