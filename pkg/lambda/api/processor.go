package api

// ProcessorInterface defines an interface for processing API Gateway proxy requests
type ProcessorInterface interface {
	NewRequestBody() interface{}
	Process(*Input) (*Output, error)
}
