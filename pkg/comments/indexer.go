package comments

import (
	"fmt"

	"github.com/bitterbattles/api/pkg/index"
)

const keyPattern = "commentIds:forBattle:%s"

// Indexer is used to index Battles in a variety of ways
type Indexer struct {
	repository index.RepositoryInterface
}

// NewIndexer creates a new indexer instance
func NewIndexer(repository index.RepositoryInterface) *Indexer {
	return &Indexer{repository}
}

// Add adds a new Comment to the relevant index(es)
func (indexer *Indexer) Add(comment *Comment) error {
	var key string
	var err error
	key = fmt.Sprintf(keyPattern, comment.BattleID)
	member := comment.ID
	score := float64(comment.CreatedOn)
	err = indexer.repository.SetScore(key, member, score)
	if err != nil {
		return err
	}
	return nil
}

// GetByBattleID gets a range of Comment IDs for the given Battle ID
func (indexer *Indexer) GetByBattleID(battleID string, page int, pageSize int) ([]string, error) {
	key := fmt.Sprintf(keyPattern, battleID)
	offset, limit := indexer.toRange(page, pageSize)
	return indexer.repository.GetRange(key, offset, limit)
}

// Delete removes the given Comment from all index(es)
func (indexer *Indexer) Delete(comment *Comment) error {
	var key string
	var err error
	key = fmt.Sprintf(keyPattern, comment.BattleID)
	member := comment.ID
	err = indexer.repository.Remove(key, member)
	if err != nil {
		return err
	}
	return nil
}

func (indexer *Indexer) toRange(page int, pageSize int) (int, int) {
	return (page - 1) * pageSize, pageSize
}
