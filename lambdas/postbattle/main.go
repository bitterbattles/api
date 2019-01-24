package postbattle

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/battles"
	"github.com/bitterbattles/api/common/bootstrap"
	"github.com/bitterbattles/api/common/loggers"
)

func main() {
	session := bootstrap.NewSession()
	dynamoClient := bootstrap.NewDynamoClient(session)
	logger := loggers.NewCloudWatchLogger()
	repository := battles.NewRepository(dynamoClient)
	handler := NewHandler(repository, logger)
	lambda.StartHandler(handler)
}
