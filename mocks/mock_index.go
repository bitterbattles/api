package mocks

// MockIndex is a mock implementation of battles.IndexInterface
type MockIndex struct {
	data map[string]*[]string
}

// NewMockIndex creates a new MockIndex instance
func NewMockIndex() *MockIndex {
	return &MockIndex{make(map[string]*[]string)}
}

// Add adds a value to the given index
func (index *MockIndex) Add(indexName string, value string) error {
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
func (index *MockIndex) GetRange(indexName string, offset int, limit int) ([]string, error) {
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
func (index *MockIndex) Reset() {
	index.data = make(map[string]*[]string)
}
