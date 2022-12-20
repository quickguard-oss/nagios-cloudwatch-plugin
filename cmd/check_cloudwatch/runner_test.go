package main

import (
	goerrors "errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	awstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/alert"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/cloudwatch/client"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/cloudwatch/client/types"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/test/helper"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/test/mock"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func Test_run(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		commandArgs             []string
		cloudwatchClientFactory func() (types.Client, error)
	}

	type testCase struct {
		name     string
		args     args
		expected alert.ReturnCode
	}

	testCases := []testCase{
		{
			name: "ok",
			args: args{
				commandArgs: []string{
					"--warning",
					"0.0:1.5",
					"--critical",
					"0.0:2.5",
					"--datapoints",
					"1/1",
					"--queries",
					`[{"Id":"e1","Expression":"TIME_SERIES(1)"}]`,
				},
				cloudwatchClientFactory: func() (types.Client, error) {
					m := &mock.CloudWatchClient{}

					output := &cloudwatch.GetMetricDataOutput{
						MetricDataResults: []awstypes.MetricDataResult{
							{
								Id: aws.String("e1"),
								Timestamps: []time.Time{
									time.Date(2022, time.September, 19, 10, 20, 30, 0, time.UTC),
								},
								Values: []float64{
									1.0,
								},
							},
						},
					}

					m.On("GetMetricData", testifymock.Anything, testifymock.Anything).Return(output, nil)

					return m, nil
				},
			},
			expected: alert.OK,
		},
		{
			name: "warning",
			args: args{
				commandArgs: []string{
					"--warning",
					"0.0:0.5",
					"--critical",
					"0.0:2.5",
					"--datapoints",
					"1/1",
					"--queries",
					`[{"Id":"e1","Expression":"TIME_SERIES(1)"}]`,
				},
				cloudwatchClientFactory: func() (types.Client, error) {
					m := &mock.CloudWatchClient{}

					output := &cloudwatch.GetMetricDataOutput{
						MetricDataResults: []awstypes.MetricDataResult{
							{
								Id: aws.String("e1"),
								Timestamps: []time.Time{
									time.Date(2022, time.September, 19, 10, 20, 30, 0, time.UTC),
								},
								Values: []float64{
									1.0,
								},
							},
						},
					}

					m.On("GetMetricData", testifymock.Anything, testifymock.Anything).Return(output, nil)

					return m, nil
				},
			},
			expected: alert.Warning,
		},
		{
			name: "critical",
			args: args{
				commandArgs: []string{
					"--warning",
					"0.0:1.5",
					"--critical",
					"0.0:0.5",
					"--datapoints",
					"1/1",
					"--queries",
					`[{"Id":"e1","Expression":"TIME_SERIES(1)"}]`,
				},
				cloudwatchClientFactory: func() (types.Client, error) {
					m := &mock.CloudWatchClient{}

					output := &cloudwatch.GetMetricDataOutput{
						MetricDataResults: []awstypes.MetricDataResult{
							{
								Id: aws.String("e1"),
								Timestamps: []time.Time{
									time.Date(2022, time.September, 19, 10, 20, 30, 0, time.UTC),
								},
								Values: []float64{
									1.0,
								},
							},
						},
					}

					m.On("GetMetricData", testifymock.Anything, testifymock.Anything).Return(output, nil)

					return m, nil
				},
			},
			expected: alert.Critical,
		},
		{
			name: "invalid args",
			args: args{
				commandArgs: []string{
					"--UNKNOWN-ARG",
				},
				cloudwatchClientFactory: client.New,
			},
			expected: alert.Unknown,
		},
		{
			name: "checker error",
			args: args{
				commandArgs: []string{
					"--warning",
					"a",
					"--critical",
					"0.0:2.5",
					"--datapoints",
					"1/1",
					"--queries",
					`[{"Id":"e1","Expression":"TIME_SERIES(1)"}]`,
				},
				cloudwatchClientFactory: client.New,
			},
			expected: alert.Unknown,
		},
		{
			name: "client error",
			args: args{
				commandArgs: []string{
					"--warning",
					"0.0:1.5",
					"--critical",
					"0.0:2.5",
					"--datapoints",
					"1/1",
					"--queries",
					"{",
				},
				cloudwatchClientFactory: client.New,
			},
			expected: alert.Unknown,
		},
		{
			name: "API error",
			args: args{
				commandArgs: []string{
					"--warning",
					"0.0:1.5",
					"--critical",
					"0.0:2.5",
					"--datapoints",
					"1/1",
					"--queries",
					`[{"Id":"e1","Expression":"TIME_SERIES(1)"}]`,
				},
				cloudwatchClientFactory: func() (types.Client, error) {
					m := &mock.CloudWatchClient{}

					output := &cloudwatch.GetMetricDataOutput{}

					m.On("GetMetricData", testifymock.Anything, testifymock.Anything).Return(output, goerrors.New(""))

					return m, nil
				},
			},
			expected: alert.Unknown,
		},
		{
			name: "evaluation error",
			args: args{
				commandArgs: []string{
					"--warning",
					"0.0:1.5",
					"--critical",
					"0.0:2.5",
					"--datapoints",
					"1/1",
					"--queries",
					`[{"Id":"e1","Expression":"TIME_SERIES(1)"}]`,
				},
				cloudwatchClientFactory: func() (types.Client, error) {
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
			},
			expected: alert.Unknown,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			helper.SetLogOutputDiscard(t)
			helper.SetCommandArgs(t, tc.args.commandArgs)
			helper.SetCloudWatchClientFactory(t, tc.args.cloudwatchClientFactory)

			assert.Equal(tc.expected, run(), "alert.ReturnCode")
		})
	}
}
