package battlesstream

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/battles"
	"github.com/bitterbattles/api/common/bootstrap"
)

func main() {
	redisClient := bootstrap.NewRedisClient()
	index := battles.NewIndex(redisClient)
	handler := NewHandler(index)
	lambda.StartHandler(handler)
}
