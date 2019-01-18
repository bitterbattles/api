package main

import (
	"github.com/bitterbattles/api/battles"
	"github.com/bitterbattles/api/core/loggers"
	"github.com/bitterbattles/api/lambda"
)

func main() {
	logger := loggers.NewCloudWatchLogger()
	index := battles.NewIndex()
	repository := battles.NewRepository()
	manager := battles.NewManager(index, repository)
	controller := battles.NewController(manager)
	handler := lambda.NewHandler(controller.HandlePost, logger)
	handler.Run()
}
