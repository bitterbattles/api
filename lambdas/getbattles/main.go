package getbattles

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/battles"
	"github.com/bitterbattles/api/common/bootstrap"
)

func main() {
	session := bootstrap.NewSession()
	dynamoClient := bootstrap.NewDynamoClient(session)
	redisClient := bootstrap.NewRedisClient()
	repository := battles.NewRepository(dynamoClient)
	index := battles.NewIndex(redisClient)
	handler := NewHandler(index, repository)
	lambda.StartHandler(handler)
}
