package battles

import "strings"

const defaultSort string = "recent"
const defaultPageSize int = 50
const maxPageSize int = 100

// ManagerInterface defines an interface for a Battle manager
type ManagerInterface interface {
	Create(string, string) (*Battle, error)
	GetPage(string, int, int) ([]*Battle, error)
}

// Manager is used to perform business logic related to Battles
type Manager struct {
	index IndexInterface
	table TableInterface
}

// NewManager creates a new Manager instance
func NewManager(index IndexInterface, table TableInterface) *Manager {
	return &Manager{index, table}
}

// Create creates a new Battle
func (manager *Manager) Create(title string, description string) (*Battle, error) {
	// TODO: Input validation
	battle := Battle{
		ID:          "TODO",
		UserID:      "TODO",
		Title:       title,
		Description: description,
		CreatedOn:   0, // TODO
	}
	err := manager.table.Add(battle)
	if err != nil {
		return nil, err
	}
	return &battle, nil
}

// GetPage gets a page of Battles
func (manager *Manager) GetPage(sort string, page int, pageSize int) ([]*Battle, error) {
	ids, err := manager.getIdsPage(sort, page, pageSize)
	if err != nil {
		return nil, err
	}
	results := make([]*Battle, 0, len(ids))
	for _, id := range ids {
		result, err := manager.getByID(id)
		if err != nil {
			return nil, err
		}
		if result != nil {
			results = append(results, result) // TODO: Make this more efficient
		}
	}
	return results, nil
}

func (manager *Manager) getIdsPage(sort string, page int, pageSize int) ([]string, error) {
	sort = manager.sanitizeSort(sort)
	page, pageSize = manager.sanitizePagination(page, pageSize)
	offset := (page - 1) * pageSize
	limit := pageSize
	results, err := manager.index.GetRange(sort, offset, limit)
	if err != nil {
		return nil, err
	}
	return []string(results), nil
}

func (manager *Manager) getByID(id string) (*Battle, error) {
	battle, err := manager.table.GetByID(id)
	if err != nil {
		return nil, err
	}
	return battle, nil
}

func (manager *Manager) sanitizePagination(page int, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	} else if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}

func (manager *Manager) sanitizeSort(sort string) string {
	sort = strings.ToLower(strings.TrimSpace(sort))
	switch sort {
	case defaultSort:
	case "popular":
	case "controversial":
		return sort
	}
	return defaultSort
}
