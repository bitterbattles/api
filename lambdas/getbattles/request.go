package getbattles

// Request represents a request
type Request struct {
	Sort     string `json:"sort"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}
