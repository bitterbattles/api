package index

import (
	"github.com/go-redis/redis"
)

// RepositoryInterface defines an interface for an index repository
type RepositoryInterface interface {
	SetScore(string, string, float64) error
	GetRange(string, int, int) ([]string, error)
	Remove(string, string) error
	HasMember(string, string) (bool, error)
}

// Repository is an implementation of RepositoryInterface using Redis
type Repository struct {
	client *redis.Client
}

// NewRepository creates a new Repository instance
func NewRepository(client *redis.Client) *Repository {
	return &Repository{client}
}

// SetScore sets a score for a member
func (repository *Repository) SetScore(key string, member string, score float64) error {
	entry := redis.Z{
		Member: member,
		Score:  score,
	}
	_, err := repository.client.ZAdd(key, entry).Result()
	return err
}

// GetRange gets a range of members
func (repository *Repository) GetRange(key string, offset int, limit int) ([]string, error) {
	results, err := repository.client.ZRevRangeByScore(key, redis.ZRangeBy{
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

// Remove removes a member
func (repository *Repository) Remove(key string, member string) error {
	_, err := repository.client.ZRem(key, member).Result()
	return err
}

// HasMember determines if the member exists
func (repository *Repository) HasMember(key string, member string) (bool, error) {
	_, err := repository.client.ZScore(key, member).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}
	hasMember := (err == nil)
	return hasMember, nil
}
