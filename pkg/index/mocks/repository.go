package mocks

// Entry represents an index entry
type Entry struct {
	Member string
	Score  float64
}

// Repository is a mocked implementation of index.Repository
type Repository struct {
	data      map[string][]*Entry
	lastAdded map[string]*Entry
}

// NewRepository creates a new repository instance
func NewRepository() *Repository {
	return &Repository{make(map[string][]*Entry), make(map[string]*Entry)}
}

// SetScore sets a score for a member
func (repository *Repository) SetScore(key string, member string, score float64) error {
	entry := &Entry{
		Member: member,
		Score:  score,
	}
	entries := repository.data[key]
	if entries != nil {
		for _, entry := range entries {
			if entry.Member == member {
				entry.Score = score
				return nil
			}
		}
		repository.data[key] = append(entries, entry)
	} else {
		entries = make([]*Entry, 1)
		entries[0] = entry
		repository.data[key] = entries
	}
	repository.lastAdded[key] = entry
	return nil
}

// GetRange gets a range of members
func (repository *Repository) GetRange(key string, offset int, limit int) ([]string, error) {
	_, ok := repository.data[key]
	if !ok {
		return []string{}, nil
	}
	allEntries := repository.data[key]
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
		values[i] = entries[i].Member
	}
	return values, nil
}

// Remove removes a member
func (repository *Repository) Remove(key string, member string) error {
	entries := repository.data[key]
	if entries != nil {
		for i := 0; i < len(entries); i++ {
			entry := entries[i]
			if entry.Member == member {
				entries = append(entries[:i], entries[i+1:]...)
				repository.data[key] = entries
			}
		}
	}
	return nil
}

// GetScore gets the score for the given member
func (repository *Repository) GetScore(key string, member string) float64 {
	entries := repository.data[key]
	if entries != nil {
		for _, entry := range entries {
			if entry.Member == member {
				return entry.Score
			}
		}
	}
	return 0
}

// GetLastAdded gets the most recently added entry
func (repository *Repository) GetLastAdded(key string) *Entry {
	return repository.lastAdded[key]
}

// Reset removes all data from the repository
func (repository *Repository) Reset() {
	repository.data = make(map[string][]*Entry)
	repository.lastAdded = make(map[string]*Entry)
}
