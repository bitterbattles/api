package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/pkg/common/bootstrap"
	"github.com/bitterbattles/api/pkg/ranks"
)

func main() {
	redisClient := bootstrap.NewRedisClient()
	repository := ranks.NewRepository(redisClient)
	handler := NewHandler(repository)
	lambda.StartHandler(handler)
}
