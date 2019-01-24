package battles

import (
	"fmt"

	"github.com/go-redis/redis"
)

const indexKeyPattern = "battles:%s"

// IndexInterface defines an interface for a Battle index
type IndexInterface interface {
	GetRange(string, int, int) ([]string, error)
	Set(string, string, float64) error
}

// Index is an implementation of IndexInterface that uses Redis
type Index struct {
	client *redis.Client
}

// NewIndex creates a new Index instance
func NewIndex(client *redis.Client) *Index {
	return &Index{client}
}

// GetRange gets a range of values from the index
func (index *Index) GetRange(indexName string, offset int, limit int) ([]string, error) {
	key := fmt.Sprintf(indexKeyPattern, indexName)
	results, err := index.client.ZRangeByScore(key, redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: int64(offset),
		Count:  int64(limit),
	}).Result()
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Set is used to set the score value for a Battle in the given index
func (index *Index) Set(indexName string, battleID string, score float64) error {
	key := fmt.Sprintf(indexKeyPattern, indexName)
	entry := redis.Z{
		Member: battleID,
		Score:  score,
	}
	_, err := index.client.ZAdd(key, entry).Result()
	return err
}
