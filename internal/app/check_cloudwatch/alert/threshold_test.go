package alert

import (
	"math"
	"testing"

	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/errors"
	"github.com/stretchr/testify/assert"
)

func Test_newThresholdRange(t *testing.T) {
	assert := assert.New(t)

	type expected struct {
		thresholdRange thresholdRange
		err            error
	}

	type testCase struct {
		name     string
		args     string
		expected expected
	}

	testCases := []testCase{
		{
			name: "end only",
			args: "0.1",
			expected: expected{
				thresholdRange: thresholdRange{
					enable:  true,
					start:   0.0,
					end:     0.1,
					inverse: false,
				},
				err: nil,
			},
		},
		{
			name: "start only",
			args: "0.2:",
			expected: expected{
				thresholdRange: thresholdRange{
					enable:  true,
					start:   0.2,
					end:     math.Inf(1),
					inverse: false,
				},
				err: nil,
			},
		},
		{
			name: "negative infinity",
			args: "~:0.3",
			expected: expected{
				thresholdRange: thresholdRange{
					enable:  true,
					start:   math.Inf(-1),
					end:     0.3,
					inverse: false,
				},
				err: nil,
			},
		},
		{
			name: "start and end",
			args: "0.4:0.5",
			expected: expected{
				thresholdRange: thresholdRange{
					enable:  true,
					start:   0.4,
					end:     0.5,
					inverse: false,
				},
				err: nil,
			},
		},
		{
			name: "negative values",
			args: "-0.7:-0.6",
			expected: expected{
				thresholdRange: thresholdRange{
					enable:  true,
					start:   -0.7,
					end:     -0.6,
					inverse: false,
				},
				err: nil,
			},
		},
		{
			name: "inverse",
			args: "@0.8:0.9",
			expected: expected{
				thresholdRange: thresholdRange{
					enable:  true,
					start:   0.8,
					end:     0.9,
					inverse: true,
				},
				err: nil,
			},
		},
		{
			name: "not specified",
			args: "",
			expected: expected{
				thresholdRange: thresholdRange{
					enable: false,
				},
				err: nil,
			},
		},
		{
			name: "illegal format",
			args: "a",
			expected: expected{
				thresholdRange: thresholdRange{},
				err:            &errors.ArgumentError{},
			},
		},
		{
			name: "start > end",
			args: "0.9:0.8",
			expected: expected{
				thresholdRange: thresholdRange{},
				err:            &errors.ArgumentError{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tr, err := newThresholdRange(tc.args)

			assert.Equal(tc.expected.thresholdRange, tr, "thresholdRange")

			if tc.expected.err != nil {
				assert.ErrorAs(err, tc.expected.err, "is error")
			} else {
				assert.Nil(err, "is not error")
			}
		})
	}
}

func Test_parseDatapointsThreshold(t *testing.T) {
	assert := assert.New(t)

	type expected struct {
		datapointsToAlarm int
		evaluationPeriods int
		err               error
	}

	type testCase struct {
		name     string
		args     string
		expected expected
	}

	testCases := []testCase{
		{
			name: "datapointsToAlarm <= evaluationPeriods",
			args: "1/1",
			expected: expected{
				datapointsToAlarm: 1,
				evaluationPeriods: 1,
				err:               nil,
			},
		},
		{
			name: "not specified",
			args: "",
			expected: expected{
				err: &errors.ArgumentError{},
			},
		},
		{
			name: "float value",
			args: "0.1/0.2",
			expected: expected{
				err: &errors.ArgumentError{},
			},
		},
		{
			name: "negative value",
			args: "-1/2",
			expected: expected{
				err: &errors.ArgumentError{},
			},
		},
		{
			name: "evaluationPeriods < datapointsToAlarm",
			args: "2/1",
			expected: expected{
				err: &errors.ArgumentError{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			datapointsToAlarm, evaluationPeriods, err := parseDatapointsThreshold(tc.args)

			assert.Equal(tc.expected.datapointsToAlarm, datapointsToAlarm, "datapointsToAlarm")
			assert.Equal(tc.expected.evaluationPeriods, evaluationPeriods, "evaluationPeriods")

			if tc.expected.err != nil {
				assert.ErrorAs(err, tc.expected.err, "is error")
			} else {
				assert.Nil(err, "is not error")
			}
		})
	}
}
