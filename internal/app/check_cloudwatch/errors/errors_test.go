package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ArgumentError_Error(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		msg   string
		key   string
		value string
	}

	type testCase struct {
		name     string
		args     args
		expected string
	}

	testCases := []testCase{
		{
			name: "error",
			args: args{
				msg:   "a",
				key:   "k",
				value: "v",
			},
			expected: `invalid argument "v" for k: a`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := NewArgumentErrorWithMessage(tc.args.msg, tc.args.key, tc.args.value)

			assert.Equal(tc.expected, err.Error(), "Error")
		})
	}
}

func Test_CloudWatchError_Error(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		name     string
		args     string
		expected string
	}

	testCases := []testCase{
		{
			name:     "error",
			args:     "a",
			expected: `unable to get metrics: a`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := NewCloudWatchError(errors.New(tc.args))

			assert.Equal(tc.expected, err.Error(), "Error")
		})
	}
}
