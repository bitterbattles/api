package loginspost

// Response represents a response body
type Response struct {
	AccessToken      string `json:"accessToken"`
	AccessExpiresIn  int    `json:"accessExpiresIn"`
	RefreshToken     string `json:"refreshToken"`
	RefreshExpiresIn int    `json:"refreshExpiresIn"`
}
