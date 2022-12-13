package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/stretchr/testify/mock"
)

type CloudWatchClient struct {
	mock.Mock
}

func (m *CloudWatchClient) GetMetricData(
	ctx context.Context,
	params *cloudwatch.GetMetricDataInput,
	optFns ...func(*cloudwatch.Options),
) (*cloudwatch.GetMetricDataOutput, error) {
	args := m.Called(ctx, params)

	return args.Get(0).(*cloudwatch.GetMetricDataOutput), args.Error(1)
}
