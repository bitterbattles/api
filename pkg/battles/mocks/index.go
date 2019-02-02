package mocks

type entry struct {
	value string
	score float64
}

// Index is a mocked implementation of battles.IndexInterface
type Index struct {
	data map[string][]*entry
}

// NewIndex creates a new Index instance
func NewIndex() *Index {
	return &Index{make(map[string][]*entry)}
}

// GetRange gets a page of values from the index
func (index *Index) GetRange(indexName string, offset int, limit int) ([]string, error) {
	_, ok := index.data[indexName]
	if !ok {
		return []string{}, nil
	}
	allEntries := index.data[indexName]
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
		values[i] = entries[i].value
	}
	return values, nil
}

// GetScore gets the score for the given value and index
func (index *Index) GetScore(indexName string, value string) float64 {
	entries := index.data[indexName]
	if entries != nil {
		for _, entry := range entries {
			if entry.value == value {
				return entry.score
			}
		}
	}
	return 0
}

// Set adds a value to the given index
func (index *Index) Set(indexName string, value string, score float64) error {
	entries := index.data[indexName]
	if entries != nil {
		for _, entry := range entries {
			if entry.value == value {
				entry.score = score
				return nil
			}
		}
		index.data[indexName] = append(entries, &entry{value, score})
	} else {
		entries = make([]*entry, 1)
		entries[0] = &entry{value, score}
		index.data[indexName] = entries
	}
	return nil
}

// Reset removes all data from the index
func (index *Index) Reset() {
	index.data = make(map[string][]*entry)
}
