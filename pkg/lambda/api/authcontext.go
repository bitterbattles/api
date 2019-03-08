package api

// AuthContext represents auth context for a request
type AuthContext struct {
	UserID    string `json:"sub"`
	CreatedOn int64  `json:"iat"`
}
