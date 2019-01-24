package loggers

// CloudWatchLogger is an implementation of LoggerInterface that uses CloudWatch
type CloudWatchLogger struct {
}

// NewCloudWatchLogger creates a new CloudWatchLogger instance
func NewCloudWatchLogger() *CloudWatchLogger {
	return &CloudWatchLogger{}
}

// Error logs an error
func (logger *CloudWatchLogger) Error(message string, err error) {
	// TODO
}
