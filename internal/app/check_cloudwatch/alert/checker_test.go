package alert

import (
	"testing"

	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/errors"
	"github.com/stretchr/testify/assert"
)

func Test_Checker_CheckStatus(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		warnRange           string
		criticalRange       string
		datapointsThreshold string
		values              []float64
	}

	type expected struct {
		returnCode ReturnCode
		err        error
	}

	type testCase struct {
		name     string
		args     args
		expected expected
	}

	testCases := []testCase{
		{
			name: "ok",
			args: args{
				warnRange:           "1.0:2.0",
				criticalRange:       "0.0:3.0",
				datapointsThreshold: "2/3",
				values: []float64{
					-0.5,
					1.5,
					1.5,
					3.5,
				},
			},
			expected: expected{
				returnCode: OK,
				err:        nil,
			},
		},
		{
			name: "warning",
			args: args{
				warnRange:           "1.0:2.0",
				criticalRange:       "0.0:3.0",
				datapointsThreshold: "2/3",
				values: []float64{
					-0.5,
					0.5,
					1.5,
					3.5,
				},
			},
			expected: expected{
				returnCode: Warning,
				err:        nil,
			},
		},
		{
			name: "warning & critical",
			args: args{
				warnRange:           "1.0:2.0",
				criticalRange:       "0.0:3.0",
				datapointsThreshold: "2/3",
				values: []float64{
					-0.5,
					-0.5,
					1.5,
					3.5,
				},
			},
			expected: expected{
				returnCode: Critical,
				err:        nil,
			},
		},
		{
			name: "critical",
			args: args{
				warnRange:           "0.0:2.0",
				criticalRange:       "1.0:3.0",
				datapointsThreshold: "2/3",
				values: []float64{
					-0.5,
					0.5,
					1.5,
					3.5,
				},
			},
			expected: expected{
				returnCode: Critical,
				err:        nil,
			},
		},
		{
			name: "insufficient data points",
			args: args{
				warnRange:           "1.0:2.0",
				criticalRange:       "0.0:3.0",
				datapointsThreshold: "2/3",
				values: []float64{
					0.5,
					1.5,
				},
			},
			expected: expected{
				returnCode: Unknown,
				err:        &errors.ArgumentError{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewChecker(tc.args.warnRange, tc.args.criticalRange, tc.args.datapointsThreshold)

			if err != nil {
				t.Error(err)
			}

			r, err := c.CheckStatus(tc.args.values)

			assert.Equal(tc.expected.returnCode, r, "ReturnCode")

			if tc.expected.err != nil {
				assert.ErrorAs(err, tc.expected.err, "is error")
			} else {
				assert.Nil(err, "is not error")
			}
		})
	}
}

func Test_Checker_Result(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		warnRange           string
		criticalRange       string
		datapointsThreshold string
		values              []float64
	}
	type expected struct {
		isWarn             bool
		isCritical         bool
		outOfWarnRange     int
		outOfCriticalRange int
	}

	type testCase struct {
		name     string
		args     args
		expected expected
	}

	testCases := []testCase{
		{
			name: "within threshold",
			args: args{
				warnRange:           "1.0:2.0",
				criticalRange:       "0.0:3.0",
				datapointsThreshold: "2/3",
				values: []float64{
					1.5,
					1.5,
					1.5,
				},
			},
			expected: expected{
				isWarn:             false,
				isCritical:         false,
				outOfWarnRange:     0,
				outOfCriticalRange: 0,
			},
		},
		{
			name: "warning count++",
			args: args{
				warnRange:           "1.0:2.0",
				criticalRange:       "0.0:3.0",
				datapointsThreshold: "2/3",
				values: []float64{
					0.5,
					1.5,
					1.5,
				},
			},
			expected: expected{
				isWarn:             false,
				isCritical:         false,
				outOfWarnRange:     1,
				outOfCriticalRange: 0,
			},
		},
		{
			name: "critical count++",
			args: args{
				warnRange:           "0.0:2.0",
				criticalRange:       "1.0:3.0",
				datapointsThreshold: "2/3",
				values: []float64{
					0.5,
					1.5,
					1.5,
				},
			},
			expected: expected{
				isWarn:             false,
				isCritical:         false,
				outOfWarnRange:     0,
				outOfCriticalRange: 1,
			},
		},
		{
			name: "is warning",
			args: args{
				warnRange:           "1.0:2.0",
				criticalRange:       "0.0:3.0",
				datapointsThreshold: "2/3",
				values: []float64{
					-0.5,
					0.5,
					1.5,
				},
			},
			expected: expected{
				isWarn:             true,
				isCritical:         false,
				outOfWarnRange:     2,
				outOfCriticalRange: 1,
			},
		},
		{
			name: "is critical",
			args: args{
				warnRange:           "0.0:2.0",
				criticalRange:       "1.0:3.0",
				datapointsThreshold: "2/3",
				values: []float64{
					-0.5,
					0.5,
					1.5,
				},
			},
			expected: expected{
				isWarn:             false,
				isCritical:         true,
				outOfWarnRange:     1,
				outOfCriticalRange: 2,
			},
		},
		{
			name: "warning & critical",
			args: args{
				warnRange:           "1.0:2.0",
				criticalRange:       "0.0:3.0",
				datapointsThreshold: "2/3",
				values: []float64{
					-0.5,
					-0.5,
					-0.5,
				},
			},
			expected: expected{
				isWarn:             false,
				isCritical:         true,
				outOfWarnRange:     3,
				outOfCriticalRange: 3,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewChecker(tc.args.warnRange, tc.args.criticalRange, tc.args.datapointsThreshold)

			if err != nil {
				t.Error(err)
			}

			c.CheckStatus(tc.args.values)

			isWarn, isCritical, outOfWarnRange, outOfCriticalRange := c.Result()

			assert.Equal(tc.expected.isWarn, isWarn, "isWarn")
			assert.Equal(tc.expected.isCritical, isCritical, "isCritical")
			assert.Equal(tc.expected.outOfWarnRange, outOfWarnRange, "outOfWarnRange")
			assert.Equal(tc.expected.outOfCriticalRange, outOfCriticalRange, "outOfCriticalRange")
		})
	}
}
