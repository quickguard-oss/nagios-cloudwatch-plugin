package cloudwatch

import (
	"context"
	goerrors "errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	awstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/cloudwatch/client"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/cloudwatch/client/types"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/errors"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/test/helper"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/test/mock"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func Test_New(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		factory    func() (types.Client, error)
		duration   int
		queriesStr string
		timeout    int
	}

	type expected struct {
		cloudWatch CloudWatch
		err        error
	}

	type testCase struct {
		name     string
		args     args
		expected expected
	}

	testCases := []testCase{
		{
			name: "success",
			args: args{
				factory:    client.New,
				duration:   10,
				queriesStr: `[{"Id":"e1","Expression":"TIME_SERIES(1)"}]`,
				timeout:    5,
			},
			expected: expected{
				cloudWatch: CloudWatch{
					duration: 10,
					queries: []awstypes.MetricDataQuery{
						{
							Id:         aws.String("e1"),
							Expression: aws.String("TIME_SERIES(1)"),
						},
					},
					timeout: 5,
				},
				err: nil,
			},
		},
		{
			name: "client error",
			args: args{
				factory: func() (types.Client, error) {
					return &mock.CloudWatchClient{}, errors.CloudWatchError{}
				},
				duration:   10,
				queriesStr: `[{"Id":"e1","Expression":"TIME_SERIES(1)"}]`,
				timeout:    5,
			},
			expected: expected{
				cloudWatch: CloudWatch{},
				err:        &errors.CloudWatchError{},
			},
		},
		{
			name: "illegal queries",
			args: args{
				factory:    client.New,
				duration:   10,
				queriesStr: `[{"a":true}]`,
				timeout:    5,
			},
			expected: expected{
				cloudWatch: CloudWatch{
					duration: 10,
					queries: []awstypes.MetricDataQuery{
						{},
					},
					timeout: 5,
				},
				err: nil,
			},
		},
		{
			name: "illegal json",
			args: args{
				factory:    client.New,
				duration:   10,
				queriesStr: "{",
				timeout:    5,
			},
			expected: expected{
				cloudWatch: CloudWatch{},
				err:        &errors.ArgumentError{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			helper.SetCloudWatchClientFactory(t, tc.args.factory)

			c, err := New(tc.args.duration, tc.args.queriesStr, tc.args.timeout)

			if tc.expected.err != nil {
				assert.ErrorAs(err, tc.expected.err, "is error")
			} else {
				assert.Nil(err, "is not error")

				assert.Equal(tc.expected.cloudWatch.duration, c.duration, "duration")
				assert.Equal(tc.expected.cloudWatch.queries, c.queries, "queries")
				assert.Equal(tc.expected.cloudWatch.timeout, c.timeout, "timeout")
			}
		})
	}
}

func Test_GetMetricValues(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		factory func() (types.Client, error)
		timeout int
	}

	type expected struct {
		values []float64
		err    error
	}

	type testCase struct {
		name     string
		args     args
		expected expected
	}

	now := time.Date(2022, time.September, 19, 10, 20, 30, 0, time.UTC)

	testCases := []testCase{
		{
			name: "success",
			args: args{
				factory: func() (types.Client, error) {
					m := &mock.CloudWatchClient{}

					input := &cloudwatch.GetMetricDataInput{
						StartTime: aws.Time(time.Date(2022, time.September, 19, 10, 10, 30, 0, time.UTC)),
						EndTime:   aws.Time(now),
						MetricDataQueries: []awstypes.MetricDataQuery{
							{
								Id:         aws.String("e1"),
								Expression: aws.String("TIME_SERIES(1)"),
							},
						},
					}

					output := &cloudwatch.GetMetricDataOutput{
						MetricDataResults: []awstypes.MetricDataResult{
							{
								Id: aws.String("e1"),
								Timestamps: []time.Time{
									now,
									time.Date(2022, time.September, 19, 10, 15, 30, 0, time.UTC),
									time.Date(2022, time.September, 19, 10, 10, 30, 0, time.UTC),
								},
								Values: []float64{
									0.0,
									0.5,
									1.0,
								},
							},
						},
					}

					m.On("GetMetricData", testifymock.Anything, input).Return(output, nil)

					return m, nil
				},
				timeout: 5,
			},
			expected: expected{
				values: []float64{
					0.0,
					0.5,
					1.0,
				},
				err: nil,
			},
		},
		{
			name: "API error",
			args: args{
				factory: func() (types.Client, error) {
					m := &mock.CloudWatchClient{}

					output := &cloudwatch.GetMetricDataOutput{}

					m.On("GetMetricData", testifymock.Anything, testifymock.Anything).Return(output, goerrors.New(""))

					return m, nil
				},
				timeout: 5,
			},
			expected: expected{
				values: []float64{},
				err:    &errors.CloudWatchError{},
			},
		},
		{
			name: "timeout",
			args: args{
				factory: func() (types.Client, error) {
					m := &mock.CloudWatchClient{}

					ctx, cancel := context.WithDeadline(
						context.Background(),
						now.Add(5*time.Second),
					)

					defer cancel()

					output := &cloudwatch.GetMetricDataOutput{}

					m.On("GetMetricData", ctx, testifymock.Anything).Return(output, goerrors.New(""))

					return m, nil
				},
				timeout: 5,
			},
			expected: expected{
				values: []float64{},
				err:    &errors.CloudWatchError{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			helper.SetCloudWatchClientFactory(t, tc.args.factory)

			c, err := New(10, `[{"Id":"e1","Expression":"TIME_SERIES(1)"}]`, tc.args.timeout)

			if err != nil {
				t.Error(err)
			}

			values, err := c.GetMetricValues(now)

			if tc.expected.err != nil {
				assert.ErrorAs(err, tc.expected.err, "is error")
			} else {
				assert.Nil(err, "is not error")

				assert.Equal(tc.expected.values, values, "values")
			}
		})
	}
}

func Test_LatestValue(t *testing.T) {
	assert := assert.New(t)

	type expected struct {
		metricName string
		value      float64
		timestamp  time.Time
	}

	type testCase struct {
		name     string
		args     func() (types.Client, error)
		expected expected
	}

	now := time.Date(2022, time.September, 19, 10, 20, 30, 0, time.UTC)

	testCases := []testCase{
		{
			name: "exists",
			args: func() (types.Client, error) {
				m := &mock.CloudWatchClient{}

				output := &cloudwatch.GetMetricDataOutput{
					MetricDataResults: []awstypes.MetricDataResult{
						{
							Id: aws.String("e1"),
							Timestamps: []time.Time{
								now,
								time.Date(2022, time.September, 19, 10, 15, 30, 0, time.UTC),
								time.Date(2022, time.September, 19, 10, 10, 30, 0, time.UTC),
							},
							Values: []float64{
								0.5,
								0.0,
								1.0,
							},
						},
					},
				}

				m.On("GetMetricData", testifymock.Anything, testifymock.Anything).Return(output, nil)

				return m, nil
			},
			expected: expected{
				metricName: "e1",
				value:      0.5,
				timestamp:  now,
			},
		},
		{
			name: "empty",
			args: func() (types.Client, error) {
				m := &mock.CloudWatchClient{}

				output := &cloudwatch.GetMetricDataOutput{
					MetricDataResults: []awstypes.MetricDataResult{
						{
							Id:         aws.String("e1"),
							Timestamps: []time.Time{},
							Values:     []float64{},
						},
					},
				}

				m.On("GetMetricData", testifymock.Anything, testifymock.Anything).Return(output, nil)

				return m, nil
			},
			expected: expected{
				metricName: "e1",
				value:      0.0,
				timestamp:  time.Time{},
			},
		},
		{
			name: "label",
			args: func() (types.Client, error) {
				m := &mock.CloudWatchClient{}

				output := &cloudwatch.GetMetricDataOutput{
					MetricDataResults: []awstypes.MetricDataResult{
						{
							Id:    aws.String("e1"),
							Label: aws.String("a"),
							Timestamps: []time.Time{
								now,
								time.Date(2022, time.September, 19, 10, 15, 30, 0, time.UTC),
								time.Date(2022, time.September, 19, 10, 10, 30, 0, time.UTC),
							},
							Values: []float64{
								0.5,
								0.0,
								1.0,
							},
						},
					},
				}

				m.On("GetMetricData", testifymock.Anything, testifymock.Anything).Return(output, nil)

				return m, nil
			},
			expected: expected{
				metricName: "a",
				value:      0.5,
				timestamp:  now,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			helper.SetCloudWatchClientFactory(t, tc.args)

			c, err := New(10, `[{"Id":"e1","Expression":"TIME_SERIES(1)"}]`, 5)

			if err != nil {
				t.Error(err)
			}

			_, err = c.GetMetricValues(now)

			if err != nil {
				t.Error(err)
			}

			metricName, value, timestamp := c.LatestValue()

			assert.Equal(tc.expected.metricName, metricName, "metricName")
			assert.Equal(tc.expected.value, value, "value")
			assert.Equal(tc.expected.timestamp, timestamp, "timestamp")
		})
	}
}
