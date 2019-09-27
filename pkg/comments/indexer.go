package comments

import (
	"fmt"

	"github.com/bitterbattles/api/pkg/index"
)

const (
	battleKeyPattern = "commentIds:forBattle:%s"
	authorKeyPattern = "commentIds:forAuthor:%s"
)

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
	member := comment.ID
	score := float64(comment.CreatedOn)
	key = fmt.Sprintf(battleKeyPattern, comment.BattleID)
	err = indexer.repository.SetScore(key, member, score)
	if err != nil {
		return err
	}
	key = fmt.Sprintf(authorKeyPattern, comment.UserID)
	err = indexer.repository.SetScore(key, member, score)
	if err != nil {
		return err
	}
	return nil
}

// GetByBattle gets a range of Comment IDs for the given Battle ID
func (indexer *Indexer) GetByBattle(battleID string, page int, pageSize int) ([]string, error) {
	key := fmt.Sprintf(battleKeyPattern, battleID)
	offset, limit := indexer.toRange(page, pageSize)
	return indexer.repository.GetRange(key, offset, limit)
}

// GetByAuthor gets a range of Comment IDs for the given User ID
func (indexer *Indexer) GetByAuthor(userID string, page int, pageSize int) ([]string, error) {
	key := fmt.Sprintf(authorKeyPattern, userID)
	offset, limit := indexer.toRange(page, pageSize)
	return indexer.repository.GetRange(key, offset, limit)
}

// IsCommentAuthor determines whether or not the given user authored the given comment
func (indexer *Indexer) IsCommentAuthor(userID string, commentID string) (bool, error) {
	key := fmt.Sprintf(authorKeyPattern, userID)
	return indexer.repository.HasMember(key, commentID)
}

// Delete removes the given Comment from all index(es)
func (indexer *Indexer) Delete(comment *Comment) error {
	var key string
	var err error
	member := comment.ID
	key = fmt.Sprintf(battleKeyPattern, comment.BattleID)
	err = indexer.repository.Remove(key, member)
	if err != nil {
		return err
	}
	key = fmt.Sprintf(authorKeyPattern, comment.UserID)
	err = indexer.repository.Remove(key, member)
	if err != nil {
		return err
	}
	return nil
}

func (indexer *Indexer) toRange(page int, pageSize int) (int, int) {
	return (page - 1) * pageSize, pageSize
}
