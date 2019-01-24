package battlesstream

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/battles"
	"github.com/bitterbattles/api/common/bootstrap"
	"github.com/bitterbattles/api/common/loggers"
)

func main() {
	session := bootstrap.NewSession()
	cloudWatchClient := bootstrap.NewCloudWatchClient(session)
	redisClient := bootstrap.NewRedisClient()
	logger := loggers.NewCloudWatchLogger(cloudWatchClient)
	index := battles.NewIndex(redisClient)
	handler := NewHandler(index, logger)
	lambda.StartHandler(handler)
}
