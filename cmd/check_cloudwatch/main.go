package main

import (
	"os"

	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/cloudwatch/client"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/container"
)

var version = "0.0.0"

func main() {
	os.Exit(
		int(run()),
	)
}

func init() {
	container.CloudWatchClientFactory = client.New
}
