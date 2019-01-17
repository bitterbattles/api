package mocks

// Index is a mocked implementation of battles.IndexInterface
type Index struct {
	data map[string]*[]string
}

// NewIndex creates a new Index instance
func NewIndex() *Index {
	return &Index{make(map[string]*[]string)}
}

// Add adds a value to the given index
func (index *Index) Add(indexName string, value string) error {
	values := index.data[indexName]
	var newValues []string
	if values == nil {
		newValues = make([]string, 1, 1)
		newValues[0] = value
	} else {
		newValues = append(*values, value)
	}
	index.data[indexName] = &newValues
	return nil
}

// GetRange gets a page of values from the index
func (index *Index) GetRange(indexName string, offset int, limit int) ([]string, error) {
	_, ok := index.data[indexName]
	if !ok {
		return []string{}, nil
	}
	values := *index.data[indexName]
	length := len(values)
	if offset >= length {
		return []string{}, nil
	}
	start := offset
	end := offset + limit
	if end > length {
		end = length
	}
	return values[start:end], nil
}

// Reset removes all data from the index
func (index *Index) Reset() {
	index.data = make(map[string]*[]string)
}
