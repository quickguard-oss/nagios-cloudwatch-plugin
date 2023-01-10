package container

import (
	"io"
	"os"

	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/cloudwatch/client/types"
)

var CloudWatchClientFactory func() (types.Client, error)
var LoggerIO io.Writer = os.Stdout

func Reset() {
	CloudWatchClientFactory = nil
	LoggerIO = os.Stdout
}

func GetCloudWatchClient() (types.Client, error) {
	return CloudWatchClientFactory()
}

func GetLoggerIO() io.Writer {
	return LoggerIO
}
