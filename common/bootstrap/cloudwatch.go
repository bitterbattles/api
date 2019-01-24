package bootstrap

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// NewCloudWatchClient creates a new CloudWatch client
func NewCloudWatchClient(session *session.Session) *cloudwatch.CloudWatch {
	return cloudwatch.New(session)
}
