package api

// Input represents parameters passed to a processor
type Input struct {
	PathParams  map[string]string
	QueryParams map[string]string
	AuthContext *AuthContext
	RequestBody interface{}
}
