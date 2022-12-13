package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/cloudwatch/client/types"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/errors"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/log"
)

func New() (types.Client, error) {
	log.V(3).Trace().
		Str("package", "cloudwatch").
		Msg("creating CloudWatch API client")

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRetryer(func() aws.Retryer {
		return aws.NopRetryer{}
	}))

	if err != nil {
		return nil, errors.NewCloudWatchError(err)
	}

	return cloudwatch.NewFromConfig(cfg), nil
}
