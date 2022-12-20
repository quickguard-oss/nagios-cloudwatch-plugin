package types

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

type Client interface {
	GetMetricData(
		ctx context.Context,
		params *cloudwatch.GetMetricDataInput,
		optFns ...func(*cloudwatch.Options),
	) (*cloudwatch.GetMetricDataOutput, error)
}
