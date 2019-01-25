package postvote

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/common/bootstrap"
	"github.com/bitterbattles/api/votes"
)

func main() {
	session := bootstrap.NewSession()
	dynamoClient := bootstrap.NewDynamoClient(session)
	repository := votes.NewRepository(dynamoClient)
	handler := NewHandler(repository)
	lambda.StartHandler(handler)
}
