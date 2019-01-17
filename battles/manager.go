package battles

import (
	"time"

	"github.com/bitterbattles/api/battles/errors"
	"github.com/bitterbattles/api/core/input"
	"github.com/rs/xid"
)

// ManagerInterface defines an interface for a Battle manager
type ManagerInterface interface {
	Create(string, string) error
	GetPage(string, int, int) ([]*Battle, error)
}

// Manager is used to perform business logic related to Battles
type Manager struct {
	index      IndexInterface
	repository RepositoryInterface
}

// NewManager creates a new Manager instance
func NewManager(index IndexInterface, repository RepositoryInterface) *Manager {
	return &Manager{index, repository}
}

// Create creates a new Battle
func (manager *Manager) Create(title string, description string) error {
	var err error
	title, err = manager.sanitizeTitle(title)
	if err != nil {
		return err
	}
	description, err = manager.sanitizeDescription(description)
	if err != nil {
		return err
	}
	battle := Battle{
		ID:          xid.New().String(),
		UserID:      "bgttr132fopt0uo06vlg",
		Title:       title,
		Description: description,
		CreatedOn:   time.Now().Unix(),
	}
	return manager.repository.Add(battle)
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
			results = append(results, result)
		}
	}
	return results, nil
}

func (manager *Manager) getIdsPage(sort string, page int, pageSize int) ([]string, error) {
	sort = manager.sanitizeSort(sort)
	page = manager.sanitizePage(page)
	pageSize = manager.sanitizePageSize(pageSize)
	offset := (page - 1) * pageSize
	limit := pageSize
	results, err := manager.index.GetRange(sort, offset, limit)
	if err != nil {
		return nil, err
	}
	return []string(results), nil
}

func (manager *Manager) getByID(id string) (*Battle, error) {
	battle, err := manager.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return battle, nil
}

func (manager *Manager) sanitizeSort(sort string) string {
	rules := input.StringRules{
		ToLower:      true,
		TrimSpace:    true,
		ValidValues:  []string{recentSort, popularSort, controversialSort},
		DefaultValue: recentSort,
	}
	sort, _ = input.SanitizeString(sort, rules, nil)
	return sort
}

func (manager *Manager) sanitizePage(page int) int {
	rules := input.IntegerRules{
		EnforceMinValue:    true,
		MinValue:           minPage,
		UseDefaultMinValue: true,
		DefaultMinValue:    defaultPage,
	}
	page, _ = input.SanitizeInteger(page, rules, nil)
	return page
}

func (manager *Manager) sanitizePageSize(pageSize int) int {
	rules := input.IntegerRules{
		EnforceMinValue:    true,
		MinValue:           minPageSize,
		UseDefaultMinValue: true,
		DefaultMinValue:    defaultPageSize,
		EnforceMaxValue:    true,
		MaxValue:           maxPageSize,
		UseDefaultMaxValue: true,
		DefaultMaxValue:    maxPageSize,
	}
	pageSize, _ = input.SanitizeInteger(pageSize, rules, nil)
	return pageSize
}

func (manager *Manager) sanitizeTitle(title string) (string, error) {
	rules := input.StringRules{
		TrimSpace: true,
		MinLength: minTitleLength,
		MaxLength: maxTitleLength,
	}
	errorCreator := func(message string) error {
		return errors.NewInvalidTitleError("Invalid title: " + message)
	}
	return input.SanitizeString(title, rules, errorCreator)
}

func (manager *Manager) sanitizeDescription(description string) (string, error) {
	rules := input.StringRules{
		TrimSpace: true,
		MinLength: minDescriptionLength,
		MaxLength: maxDescriptionLength,
	}
	errorCreator := func(message string) error {
		return errors.NewInvalidDescriptionError("Invalid description: " + message)
	}
	return input.SanitizeString(description, rules, errorCreator)
}
