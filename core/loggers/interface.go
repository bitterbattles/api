package loggers

// LoggerInterface defines an interface for loggers
type LoggerInterface interface {
	Error(string, error)
}
