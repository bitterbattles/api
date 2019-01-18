package mocks

// Logger is a mocked implementation of loggers.LoggerInterface
type Logger struct {
	lastMessage string
	lastError   error
}

// NewLogger creates a new Logger instance
func NewLogger() *Logger {
	return &Logger{}
}

// Error logs an error
func (logger *Logger) Error(message string, err error) {
	logger.lastMessage = message
	logger.lastError = err
}

// GetLastEntry gets the last message and error entry
func (logger *Logger) GetLastEntry() (string, error) {
	return logger.lastMessage, logger.lastError
}

// Reset resets last entry data
func (logger *Logger) Reset() {
	logger.lastMessage = ""
	logger.lastError = nil
}
