package battles

const tableName = "battles"

// TableInterface defines an interface for a Battles table
type TableInterface interface {
	Add(Battle) error
	GetByID(string) (*Battle, error)
}

// Table is an implementation of TableInterface DynamoDB
type Table struct {
	// client *dynamodb.DynamoDB
}

// NewTable creates a new Battles table instance
func NewTable( /*client *dynamodb.DynamoDB*/ ) *Table {
	return &Table{ /*client*/ }
}

// Add is used to insert a new Battle document
func (table *Table) Add(battle Battle) error {
	// item, err := table.toMap(battle)
	// if err != nil {
	// 	return nil, err
	// }
	// _, err = table.client.PutItem(&dynamodb.PutItemInput{
	// 	Item:      item,
	// 	TableName: aws.String(tableName),
	// })
	// return err
	return nil
}

// GetByID is used to get a Battle by ID
func (table *Table) GetByID(id string) (*Battle, error) {
	// result, err := table.client.GetItem(&dynamodb.GetItemInput{
	// 	TableName: aws.String(tableName),
	// 	Key: map[string]*dynamodb.AttributeValue{
	// 		"id": {
	// 			S: aws.String(id),
	// 		},
	// 	},
	// })
	// battle, err := table.toBattle(result)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}

// func (table *Table) toBattle(item map[string]*dynamodb.AttributeValue) (Battle, error) {
// 	battle := Battle{}
// 	err := dynamodbattribute.UnmarshalMap(item, &battle)
// 	if err != nil {
// 		return Battle{}, err
// 	}
// 	return battle, nil
// }
//
// func (table *Table) toMap(battle Battle) (map[string]*dynamodb.AttributeValue, error) {
// 	item, err := dynamodbattribute.MarshalMap(battle)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return item, nil
// }
