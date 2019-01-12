package mocks

import (
	"github.com/bitterbattles/api/battles"
)

// MockTable is a mocked implementation of battles.TableInterface
type MockTable struct {
	data map[string]*battles.Battle
}

// NewMockTable creates a new Battles table instance
func NewMockTable() *MockTable {
	return &MockTable{make(map[string]*battles.Battle)}
}

// Add is used to insert a new Battle document
func (table *MockTable) Add(battle battles.Battle) error {
	table.data[battle.ID] = &battle
	return nil
}

// GetByID is used to get a Battle by ID
func (table *MockTable) GetByID(id string) (*battles.Battle, error) {
	_, ok := table.data[id]
	if !ok {
		return nil, nil
	}
	return table.data[id], nil
}

// Reset removes all data from the table
func (table *MockTable) Reset() {
	table.data = make(map[string]*battles.Battle)
}
