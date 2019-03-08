package users

// User model
type User struct {
	ID              string `json:"id"`
	Username        string `json:"username"`        // Lower-cased username (for lookup)
	DisplayUsername string `json:"displayUsername"` // Username with original casing preserved (for display)
	PasswordHash    string `json:"passwordHash"`
	CreatedOn       int64  `json:"createdOn"`
	State           int    `json:"state"`
}
