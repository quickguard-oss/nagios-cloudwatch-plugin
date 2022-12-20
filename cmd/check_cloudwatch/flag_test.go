package main

import (
	"testing"

	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/errors"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/test/helper"
	"github.com/stretchr/testify/assert"
)

func Test_parseFlags(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		name     string
		args     []string
		expected error
	}

	testCases := []testCase{
		{
			name: "valid",
			args: []string{
				"--warning",
				"0.0:1.0",
				"--critical",
				"@~:2.0",
				"--datapoints",
				"3/4",
				"--queries",
				`{"a":true}`,
				"--duration",
				"30",
				"--timeout",
				"5",
				"--classic-output",
				"-vvv",
			},
			expected: nil,
		},
		{
			name: "unknown flag",
			args: []string{
				"--warning",
				"0.0:1.0",
				"--critical",
				"@~:2.0",
				"--datapoints",
				"3/4",
				"--queries",
				`{"a":true}`,
				"--duration",
				"30",
				"--timeout",
				"5",
				"--classic-output",
				"-vvv",
				"--UNKNOWN-ARG",
			},
			expected: &errors.ArgumentError{},
		},
		{
			name: "empty queries",
			args: []string{
				"--warning",
				"0.0:1.0",
				"--critical",
				"@~:2.0",
				"--datapoints",
				"3/4",
				"--queries",
				"",
				"--duration",
				"30",
				"--timeout",
				"5",
				"--classic-output",
				"-vvv",
			},
			expected: &errors.ArgumentError{},
		},
		{
			name: "non-positive duration",
			args: []string{
				"--warning",
				"0.0:1.0",
				"--critical",
				"@~:2.0",
				"--datapoints",
				"3/4",
				"--queries",
				`{"a":true}`,
				"--duration",
				"0",
				"--timeout",
				"5",
				"--classic-output",
				"-vvv",
			},
			expected: &errors.ArgumentError{},
		},
		{
			name: "non-positive timeout",
			args: []string{
				"--warning",
				"0.0:1.0",
				"--critical",
				"@~:2.0",
				"--datapoints",
				"3/4",
				"--queries",
				`{"a":true}`,
				"--duration",
				"30",
				"--timeout",
				"0",
				"--classic-output",
				"-vvv",
			},
			expected: &errors.ArgumentError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			helper.SetCommandArgs(t, tc.args)

			_, err := parseFlags()

			if tc.expected != nil {
				assert.ErrorAs(err, tc.expected, "is error")
			} else {
				assert.Nil(err, "is not error")
			}
		})
	}
}
