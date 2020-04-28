package comments

import (
	"fmt"

	"github.com/bitterbattles/api/pkg/elasticsearch"
)

const (
	indexName = "comments"
)

// IndexInterface defines an interface for an index
type IndexInterface interface {
	GetByBattle(string, int, int) ([]string, error)
	GetByAuthor(string, int, int) ([]string, error)
	Upsert(*Comment) error
}

// Index is used to index Comments in a variety of ways
type Index struct {
	client *elasticsearch.Client
}

// NewIndex creates a new index instance
func NewIndex(client *elasticsearch.Client) *Index {
	return &Index{
		client: client,
	}
}

// GetByBattle gets a page of Comment IDs for the given Battle ID
func (index *Index) GetByBattle(battleID string, page int, pageSize int) ([]string, error) {
	from, size := index.toRange(page, pageSize)
	body := fmt.Sprintf(`{
		"query": {
			"match": {
				"battleId": "%s"
			}
		},
		"sort": [{"createdOn": "desc"}],
		"from": %d,
		"size": %d
	}`, battleID, from, size)
	return index.client.Search(indexName, body)
}

// GetByAuthor gets a page of Comment IDs for the given User ID
func (index *Index) GetByAuthor(userID string, page int, pageSize int) ([]string, error) {
	from, size := index.toRange(page, pageSize)
	body := fmt.Sprintf(`{
		"query": {
			"match": {
				"userId": "%s"
			}
		},
		"sort": [{"createdOn": "desc"}],
		"from": %d,
		"size": %d
	}`, userID, from, size)
	return index.client.Search(indexName, body)
}

// Upsert upserts a Comment into the index
func (index *Index) Upsert(comment *Comment) error {
	body := fmt.Sprintf(`{
		"doc": {
			"battleId": "%s",
			"userId": "%s",
			"comment": "%s",
			"createdOn": %d
		},
		"doc_as_upsert": true
	}`, comment.BattleID, comment.UserID, comment.Comment, comment.CreatedOn)
	return index.client.Update(indexName, comment.ID, body)
}

func (index *Index) toRange(page int, pageSize int) (int, int) {
	return (page - 1) * pageSize, pageSize
}
