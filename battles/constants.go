package battles

// Sort constants
const (
	recentSort        = "recent"
	popularSort       = "popular"
	controversialSort = "controversial"
	defaultSort       = recentSort
)

// Pagination constants
const (
	minPage         = 1
	defaultPage     = 1
	minPageSize     = 1
	maxPageSize     = 100
	defaultPageSize = 50
)

// Battle constants
const (
	minTitleLength       = 4
	maxTitleLength       = 50
	minDescriptionLength = 4
	maxDescriptionLength = 500
	indexKeyPattern      = "battles:%s"
	tableName            = "battles"
)
