package jwt

// Header represents a JWT header
type Header struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}
