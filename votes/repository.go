package votes

// RepositoryInterface defines an interface for a Vote repository
type RepositoryInterface interface {
	Add(Vote) error
}

// Repository is an implementation of RepositoryInterface using DynamoDB
type Repository struct {
	// client *dynamodb.DynamoDB
}

// NewRepository creates a new Votes repository instance
func NewRepository( /*client *dynamodb.DynamoDB*/ ) *Repository {
	return &Repository{ /*client*/ }
}

// Add is used to insert a new Vote document
func (repository *Repository) Add(vote Vote) error {
	return nil
}
