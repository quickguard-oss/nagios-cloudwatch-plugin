package log

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_SetVerbosity(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		verbosity int
	}

	type expected []zerolog.Level

	type testCase struct {
		name     string
		args     args
		expected expected
	}

	testCases := []testCase{
		{
			name: "v0",
			args: args{
				verbosity: 0,
			},
			expected: expected{
				zerolog.TraceLevel,
				zerolog.Disabled,
				zerolog.Disabled,
				zerolog.Disabled,
			},
		},
		{
			name: "v1",
			args: args{
				verbosity: 1,
			},
			expected: expected{
				zerolog.TraceLevel,
				zerolog.TraceLevel,
				zerolog.Disabled,
				zerolog.Disabled,
			},
		},
		{
			name: "v2",
			args: args{
				verbosity: 2,
			},
			expected: expected{
				zerolog.TraceLevel,
				zerolog.TraceLevel,
				zerolog.TraceLevel,
				zerolog.Disabled,
			},
		},
		{
			name: "v3",
			args: args{
				verbosity: 3,
			},
			expected: expected{
				zerolog.TraceLevel,
				zerolog.TraceLevel,
				zerolog.TraceLevel,
				zerolog.TraceLevel,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			SetVerbosity(tc.args.verbosity)

			t.Cleanup(func() {
				Reset()
			})

			assert.Equal(tc.expected[0], V(0).GetLevel(), "v0 logger level")
			assert.Equal(tc.expected[1], V(1).GetLevel(), "v1 logger level")
			assert.Equal(tc.expected[2], V(2).GetLevel(), "v2 logger level")
			assert.Equal(tc.expected[3], V(3).GetLevel(), "v3 logger level")
		})
	}
}
