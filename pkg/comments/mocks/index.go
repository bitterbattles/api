package mocks

import (
	"github.com/bitterbattles/api/pkg/comments"
)

// Index is used to index Comments in a variety of ways
type Index struct {
	entries          []*comments.Comment
	lastCommentAdded *comments.Comment
}

// NewIndex creates a new Index instance
func NewIndex() *Index {
	return &Index{
		entries:          make([]*comments.Comment, 0),
		lastCommentAdded: nil,
	}
}

// GetByBattle gets a page of Comment IDs for the given Battle ID
func (index *Index) GetByBattle(battleID string, page int, pageSize int) ([]string, error) {
	entriesByBattle := make([]*comments.Comment, 0, len(index.entries))
	for _, comment := range index.entries {
		if comment.BattleID == battleID {
			entriesByBattle = append(entriesByBattle, comment)
		}
	}
	ids := make([]string, 0, pageSize)
	max := len(entriesByBattle)
	start := (page - 1) * pageSize
	if start >= max {
		return ids, nil
	}
	end := start + pageSize
	if end > max {
		end = max
	}
	entries := entriesByBattle[start:end]
	for _, comment := range entries {
		ids = append(ids, comment.ID)
	}
	return ids, nil
}

// GetByAuthor gets a page of Comment IDs for the given User ID
func (index *Index) GetByAuthor(userID string, page int, pageSize int) ([]string, error) {
	entriesByAuthor := make([]*comments.Comment, 0, len(index.entries))
	for _, comment := range index.entries {
		if comment.UserID == userID {
			entriesByAuthor = append(entriesByAuthor, comment)
		}
	}
	ids := make([]string, 0, pageSize)
	max := len(entriesByAuthor)
	start := (page - 1) * pageSize
	if start >= max {
		return ids, nil
	}
	end := start + pageSize
	if end > max {
		end = max
	}
	entries := entriesByAuthor[start:end]
	for _, comment := range entries {
		ids = append(ids, comment.ID)
	}
	return ids, nil
}

// Upsert upserts a Comment into the index
func (index *Index) Upsert(comment *comments.Comment) error {
	found := false
	for i, existingComment := range index.entries {
		if existingComment.ID == comment.ID {
			index.entries[i] = comment
			found = true
			break
		}
	}
	if !found {
		index.entries = append(index.entries, comment)
	}
	index.lastCommentAdded = comment
	return nil
}

// GetLastCommentAdded gets the last Comment added
func (index *Index) GetLastCommentAdded() *comments.Comment {
	return index.lastCommentAdded
}
