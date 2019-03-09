package mocks

// Rank represents a ranking within a category
type Rank struct {
	BattleID string
	Score    float64
}

// Repository is a mocked implementation of ranks.RepositoryInterface
type Repository struct {
	data      map[string][]*Rank
	lastAdded *Rank
}

// NewRepository creates a new Rank repository instance
func NewRepository() *Repository {
	return &Repository{make(map[string][]*Rank), nil}
}

// DeleteByBattleID deletes a score in a category for the given Battle ID
func (repository *Repository) DeleteByBattleID(category string, battleID string) error {
	entries := repository.data[category]
	if entries != nil {
		for i := 0; i < len(entries); i++ {
			entry := entries[i]
			if entry.BattleID == battleID {
				entries = append(entries[:i], entries[i+1:]...)
				repository.data[category] = entries
			}
		}
	}
	return nil
}

// GetRankedBattleIDs gets a range of Battle IDs by category, sorted by score
func (repository *Repository) GetRankedBattleIDs(category string, offset int, limit int) ([]string, error) {
	allEntries := repository.data[category]
	if allEntries == nil {
		return []string{}, nil
	}
	length := len(allEntries)
	if offset >= length {
		return []string{}, nil
	}
	start := offset
	end := offset + limit
	if end > length {
		end = length
	}
	entries := allEntries[start:end]
	values := make([]string, len(entries))
	for i := 0; i < len(values); i++ {
		values[i] = entries[i].BattleID
	}
	return values, nil
}

// GetScore gets the score for the given Battle and category
func (repository *Repository) GetScore(category string, battleID string) float64 {
	entries := repository.data[category]
	if entries != nil {
		for _, entry := range entries {
			if entry.BattleID == battleID {
				return entry.Score
			}
		}
	}
	return 0
}

// SetScore is used to add or update the score value for a Battle in the given ranking type
func (repository *Repository) SetScore(category string, battleID string, score float64) error {
	entry := Rank{
		BattleID: battleID,
		Score:    score,
	}
	entries := repository.data[category]
	if entries != nil {
		for _, entry := range entries {
			if entry.BattleID == battleID {
				entry.Score = score
				return nil
			}
		}
		repository.data[category] = append(entries, &entry)
	} else {
		entries = make([]*Rank, 1)
		entries[0] = &entry
		repository.data[category] = entries
	}
	return nil
}

// GetLastAdded gets the most recently added Rank
func (repository *Repository) GetLastAdded() *Rank {
	return repository.lastAdded
}

// Reset removes all data from the repository
func (repository *Repository) Reset() {
	repository.data = make(map[string][]*Rank)
	repository.lastAdded = nil
}
