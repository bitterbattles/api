package postvote

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/common/bootstrap"
	"github.com/bitterbattles/api/common/loggers"
	"github.com/bitterbattles/api/votes"
)

func main() {
	session := bootstrap.NewSession()
	cloudWatchClient := bootstrap.NewCloudWatchClient(session)
	dynamoClient := bootstrap.NewDynamoClient(session)
	logger := loggers.NewCloudWatchLogger(cloudWatchClient)
	repository := votes.NewRepository(dynamoClient)
	handler := NewHandler(repository, logger)
	lambda.StartHandler(handler)
}
