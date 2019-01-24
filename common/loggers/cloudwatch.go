package loggers

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// CloudWatchLogger is an implementation of LoggerInterface that uses CloudWatch
type CloudWatchLogger struct {
	client *cloudwatch.CloudWatch
}

// NewCloudWatchLogger creates a new CloudWatchLogger instance
func NewCloudWatchLogger(client *cloudwatch.CloudWatch) *CloudWatchLogger {
	return &CloudWatchLogger{client}
}

// Error logs an error
func (logger *CloudWatchLogger) Error(message string, err error) {
	// namespace := "TODO"
	// timestamp := time.Now()
	// metricData := make([]*cloudwatch.MetricDatum, 1)
	// metricData[0] = &cloudwatch.MetricDatum{
	// 	Timestamp: &timestamp,
	// 	// TODO
	// }
	// request, _ := logger.client.PutMetricDataRequest(&cloudwatch.PutMetricDataInput{
	// 	Namespace:  &namespace,
	// 	MetricData: metricData,
	// })
	// request.Send()
}
