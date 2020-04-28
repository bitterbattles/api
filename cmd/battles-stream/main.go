package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bitterbattles/api/pkg/battles"
	"github.com/bitterbattles/api/pkg/elasticsearch"
	"github.com/bitterbattles/api/pkg/lambda/stream"
)

func main() {
	esClient := elasticsearch.NewClient(os.Getenv("ELASTICSEARCH_ADDRESS"))
	battlesIndex := battles.NewIndex(esClient)
	battlesScorer := battles.NewScorer()
	processor := NewProcessor(battlesIndex, battlesScorer)
	handler := stream.NewHandler(processor)
	lambda.StartHandler(handler)
}
