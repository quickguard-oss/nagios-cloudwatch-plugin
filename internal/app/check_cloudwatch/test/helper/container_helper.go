package helper

import (
	"io"
	"testing"

	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/cloudwatch/client/types"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/container"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/log"
)

func SetCloudWatchClientFactory(t *testing.T, factory func() (types.Client, error)) {
	t.Helper()

	t.Cleanup(func() {
		container.Reset()
	})

	orig := container.CloudWatchClientFactory

	container.CloudWatchClientFactory = factory

	t.Cleanup(func() {
		container.CloudWatchClientFactory = orig
	})
}

func SetLoggerIO(t *testing.T, logger io.Writer) {
	t.Helper()

	t.Cleanup(func() {
		log.Reset()
	})

	loggerIO := container.LoggerIO

	container.LoggerIO = logger

	t.Cleanup(func() {
		container.LoggerIO = loggerIO
	})
}

func SetLogOutputDiscard(t *testing.T) {
	t.Helper()

	SetLoggerIO(t, io.Discard)
}
