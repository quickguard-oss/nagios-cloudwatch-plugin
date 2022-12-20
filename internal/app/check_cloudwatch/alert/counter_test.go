package alert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_counter_examine(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		thresholdRange string
		values         []float64
	}

	type testCase struct {
		name     string
		args     args
		expected int
	}

	testCases := []testCase{
		{
			name: "range",
			args: args{
				thresholdRange: "0.0:1.0",
				values: []float64{
					-0.1,
					0,
					0.1,
					0.5,
					0.9,
					1.0,
					1.1,
				},
			},
			expected: 2,
		},
		{
			name: "inverse range",
			args: args{
				thresholdRange: "@0.0:1.0",
				values: []float64{
					-0.1,
					0,
					0.1,
					0.5,
					0.9,
					1.0,
					1.1,
				},
			},
			expected: 5,
		},
		{
			name: "not specified",
			args: args{
				thresholdRange: "",
				values: []float64{
					-0.1,
					0,
					0.1,
					0.5,
					0.9,
					1.0,
					1.1,
				},
			},
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tr, err := newThresholdRange(tc.args.thresholdRange)

			if err != nil {
				t.Error(err)
			}

			c := newCounter(tr, 0)

			for _, v := range tc.args.values {
				c.examine(v)
			}

			assert.Equal(tc.expected, c.count, "count")
		})
	}
}

func Test_counter_over(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		count             int
		datapointsToAlarm int
	}

	type testCase struct {
		name     string
		args     args
		expected bool
	}

	testCases := []testCase{
		{
			name: "count < datapointsToAlarm",
			args: args{
				count:             1,
				datapointsToAlarm: 2,
			},
			expected: false,
		},
		{
			name: "count == datapointsToAlarm",
			args: args{
				count:             2,
				datapointsToAlarm: 2,
			},
			expected: true,
		},
		{
			name: "count > datapointsToAlarm",
			args: args{
				count:             3,
				datapointsToAlarm: 2,
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tr, err := newThresholdRange("")

			if err != nil {
				t.Error(err)
			}

			c := newCounter(tr, tc.args.datapointsToAlarm)

			c.count = tc.args.count

			assert.Equal(tc.expected, c.over(), "over")
		})
	}
}
