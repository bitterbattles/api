package battles

import (
	"fmt"

	"github.com/bitterbattles/api/pkg/elasticsearch"
)

const (
	indexName         = "battles"
	popularSort       = "popular"
	controversialSort = "controversial"
)

// IndexInterface defines an interface for an index
type IndexInterface interface {
	GetGlobal(string, int, int) ([]string, error)
	GetByAuthor(string, string, int, int) ([]string, error)
	Upsert(*Battle, int64, int64) error
	Delete(*Battle) error
}

// Index is used to index Battles in a variety of ways
type Index struct {
	client *elasticsearch.Client
}

// NewIndex creates a new Index instance
func NewIndex(client *elasticsearch.Client) *Index {
	return &Index{
		client: client,
	}
}

// GetGlobal gets a page of global Battle IDs
func (index *Index) GetGlobal(sort string, page int, pageSize int) ([]string, error) {
	sortField := index.toSortField(sort)
	from, size := index.toRange(page, pageSize)
	body := fmt.Sprintf(`{
		"query": {
			"match_all": {}
		},
		"sort": [{"%s": "desc"}],
		"from": %d,
		"size": %d
	}`, sortField, from, size)
	return index.client.Search(indexName, body)
}

// GetByAuthor gets a page of Battle IDs authored by the given user ID
func (index *Index) GetByAuthor(userID string, sort string, page int, pageSize int) ([]string, error) {
	sortField := index.toSortField(sort)
	from, size := index.toRange(page, pageSize)
	body := fmt.Sprintf(`{
		"query": {
			"match": {
				"userId": "%s"
			}
		},
		"sort": [{"%s": "desc"}],
		"from": %d,
		"size": %d
	}`, userID, sortField, from, size)
	return index.client.Search(indexName, body)
}

// Upsert upserts a Battle into the index
func (index *Index) Upsert(battle *Battle, popularity int64, controversy int64) error {
	body := fmt.Sprintf(`{
		"doc": {
			"userId": "%s",
			"title": "%s",
			"description": "%s",
			"votesFor": %d,
			"votesAgainst": %d,
			"comments": %d,
			"createdOn": %d,
			"popularity": %d,
			"controversy": %d
		},
		"doc_as_upsert": true
	}`, battle.UserID, battle.Title, battle.Description, battle.VotesFor, battle.VotesAgainst, battle.Comments, battle.CreatedOn, popularity, controversy)
	return index.client.Update(indexName, battle.ID, body)
}

// Delete deletes a Battle from the index
func (index *Index) Delete(battle *Battle) error {
	return index.client.Delete(indexName, battle.ID)
}

func (index *Index) toSortField(sort string) string {
	switch sort {
	case popularSort:
		return "popularity"
	case controversialSort:
		return "controversy"
	default:
		return "createdOn"
	}
}

func (index *Index) toRange(page int, pageSize int) (int, int) {
	return (page - 1) * pageSize, pageSize
}
