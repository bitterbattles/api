package battles

// IndexInterface is an interface for a Battles index
type IndexInterface interface {
	GetRange(string, int, int) ([]string, error)
}

// Index is an implementation of IndexInterface that uses Redis
type Index struct {
	// client *redis.Client
}

// NewIndex creates a new Index instance
func NewIndex( /*client *redis.Client*/ ) *Index {
	return &Index{ /*client*/ }
}

// GetRange gets a range of values from the index
func (index *Index) GetRange(indexName string, offset int, limit int) ([]string, error) {
	// key := fmt.Sprintf("indexes:%s", strings.ToLower(indexName))
	// results, err := index.client.ZRangeByScore(this.key, redis.ZRangeBy{
	// 	Min:    "-inf",
	// 	Max:    "+inf",
	// 	Offset: int64(page * pageSize),
	// 	Count:  int64(pageSize),
	// }).Result()
	// if err != nil {
	// 	return nil, err
	// }
	// return results, nil
	return nil, nil
}
