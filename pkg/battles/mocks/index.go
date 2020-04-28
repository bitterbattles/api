package mocks

import (
	"github.com/bitterbattles/api/pkg/battles"
)

// Index is used to index Battles in a variety of ways
type Index struct {
	entries              []*battles.Battle
	lastBattleAdded      *battles.Battle
	lastPopularityAdded  int64
	lastControversyAdded int64
}

// NewIndex creates a new Index instance
func NewIndex() *Index {
	return &Index{
		entries:              make([]*battles.Battle, 0),
		lastBattleAdded:      nil,
		lastPopularityAdded:  0,
		lastControversyAdded: 0,
	}
}

// GetGlobal gets a page of global Battle IDs
func (index *Index) GetGlobal(sort string, page int, pageSize int) ([]string, error) {
	ids := make([]string, 0, pageSize)
	max := len(index.entries)
	start := (page - 1) * pageSize
	if start >= max {
		return ids, nil
	}
	end := start + pageSize
	if end > max {
		end = max
	}
	entries := index.entries[start:end]
	for _, battle := range entries {
		ids = append(ids, battle.ID)
	}
	return ids, nil
}

// GetByAuthor gets a page of Battle IDs authored by the given user ID
func (index *Index) GetByAuthor(userID string, sort string, page int, pageSize int) ([]string, error) {
	entriesByAuthor := make([]*battles.Battle, 0, len(index.entries))
	for _, battle := range index.entries {
		if battle.UserID == userID {
			entriesByAuthor = append(entriesByAuthor, battle)
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
	for _, battle := range entries {
		ids = append(ids, battle.ID)
	}
	return ids, nil
}

// Upsert upserts a Battle into the index
func (index *Index) Upsert(battle *battles.Battle, popularity int64, controversy int64) error {
	found := false
	for i, existingBattle := range index.entries {
		if existingBattle.ID == battle.ID {
			index.entries[i] = battle
			found = true
			break
		}
	}
	if !found {
		index.entries = append(index.entries, battle)
	}
	index.lastBattleAdded = battle
	index.lastPopularityAdded = popularity
	index.lastControversyAdded = controversy
	return nil
}

// Delete deletes a Battle from the index
func (index *Index) Delete(battle *battles.Battle) error {
	found := -1
	for i, existingBattle := range index.entries {
		if existingBattle.ID == battle.ID {
			index.entries[i] = battle
			found = i
			break
		}
	}
	if found >= 0 {
		index.entries = append(index.entries[:found], index.entries[found+1:]...)
	}
	return nil
}

// GetLastBattleAdded gets the last Battle added
func (index *Index) GetLastBattleAdded() *battles.Battle {
	return index.lastBattleAdded
}

// GetLastPopularityAdded gets the last popularity score added
func (index *Index) GetLastPopularityAdded() int64 {
	return index.lastPopularityAdded
}

// GetLastControversyAdded gets the last controversy score added
func (index *Index) GetLastControversyAdded() int64 {
	return index.lastControversyAdded
}
