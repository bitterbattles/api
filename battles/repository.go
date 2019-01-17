package battles

// RepositoryInterface defines an interface for a Battle repository
type RepositoryInterface interface {
	Add(Battle) error
	GetByID(string) (*Battle, error)
}

// Repository is an implementation of RepositoryInterface using DynamoDB
type Repository struct {
	// client *dynamodb.DynamoDB
}

// NewRepository creates a new Battles repository instance
func NewRepository( /*client *dynamodb.DynamoDB*/ ) *Repository {
	return &Repository{ /*client*/ }
}

// Add is used to insert a new Battle document
func (repository *Repository) Add(battle Battle) error {
	// item, err := repository.toMap(battle)
	// if err != nil {
	// 	return nil, err
	// }
	// _, err = repository.client.PutItem(&dynamodb.PutItemInput{
	// 	Item:      item,
	// 	TableName: aws.String(tableName),
	// })
	// return err
	return nil
}

// GetByID is used to get a Battle by ID
func (repository *Repository) GetByID(id string) (*Battle, error) {
	// result, err := repository.client.GetItem(&dynamodb.GetItemInput{
	// 	TableName: aws.String(tableName),
	// 	Key: map[string]*dynamodb.AttributeValue{
	// 		"id": {
	// 			S: aws.String(id),
	// 		},
	// 	},
	// })
	// battle, err := repository.toBattle(result)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}

// func (repository *Repository) toBattle(item map[string]*dynamodb.AttributeValue) (Battle, error) {
// 	battle := Battle{}
// 	err := dynamodbattribute.UnmarshalMap(item, &battle)
// 	if err != nil {
// 		return Battle{}, err
// 	}
// 	return battle, nil
// }
//
// func (repository *Repository) toMap(battle Battle) (map[string]*dynamodb.AttributeValue, error) {
// 	item, err := dynamodbattribute.MarshalMap(battle)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return item, nil
// }
