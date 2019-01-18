package main

import (
	"github.com/bitterbattles/api/core/loggers"
	"github.com/bitterbattles/api/lambda"
	"github.com/bitterbattles/api/votes"
)

func main() {
	logger := loggers.NewCloudWatchLogger()
	repository := votes.NewRepository()
	manager := votes.NewManager(repository)
	controller := votes.NewController(manager)
	handler := lambda.NewHandler(controller.HandlePost, logger)
	handler.Run()
}
