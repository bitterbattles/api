package elasticsearch

// Result is the top-level result object
type Result struct {
	Hits *Hits `json:"hits"`
}

// Hits is the description of returned results
type Hits struct {
	Total *HitTotal `json:"total"`
	Hits  []*Hit    `json:"hits"`
}

// HitTotal describes the total hits
type HitTotal struct {
	Value int `json:"value"`
}

// Hit is an individual result
type Hit struct {
	ID string `json:"_id"`
}
