package getbattles

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/battles"
	"github.com/bitterbattles/api/common/bootstrap"
	"github.com/bitterbattles/api/common/loggers"
)

func main() {
	session := bootstrap.NewSession()
	dynamoClient := bootstrap.NewDynamoClient(session)
	redisClient := bootstrap.NewRedisClient()
	logger := loggers.NewCloudWatchLogger()
	repository := battles.NewRepository(dynamoClient)
	index := battles.NewIndex(redisClient)
	handler := NewHandler(index, repository, logger)
	lambda.StartHandler(handler)
}
