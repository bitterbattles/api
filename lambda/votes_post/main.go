package main

import (
	"github.com/bitterbattles/api/lambda"
	"github.com/bitterbattles/api/votes"
)

func main() {
	repository := votes.NewRepository()
	manager := votes.NewManager(repository)
	controller := votes.NewController(manager)
	handler := lambda.NewHandler(controller.HandlePost)
	handler.Run()
}
