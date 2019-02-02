package ranks

import (
	"fmt"

	"github.com/go-redis/redis"
)

const keyPattern = "ranks:%s"

// RepositoryInterface defines an interface for a Rank repository
type RepositoryInterface interface {
	GetRankedBattleIDs(string, int, int) ([]string, error)
	SetScore(string, string, float64) error
}

// Repository is an implementation of RepositoryInterface using Redis
type Repository struct {
	client *redis.Client
}

// NewRepository creates a new Rank repository instance
func NewRepository(client *redis.Client) *Repository {
	return &Repository{client}
}

// GetRankedBattleIDs gets a range of Battle IDs by category, sorted by score
func (repository *Repository) GetRankedBattleIDs(category string, offset int, limit int) ([]string, error) {
	key := fmt.Sprintf(keyPattern, category)
	results, err := repository.client.ZRangeByScore(key, redis.ZRangeBy{
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

// SetScore is used to add or update the score value for a Battle in the given ranking type
func (repository *Repository) SetScore(category string, battleID string, score float64) error {
	key := fmt.Sprintf(keyPattern, category)
	entry := redis.Z{
		Member: battleID,
		Score:  score,
	}
	_, err := repository.client.ZAdd(key, entry).Result()
	return err
}
