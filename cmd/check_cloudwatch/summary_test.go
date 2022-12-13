package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/alert"
	"github.com/stretchr/testify/assert"
)

func captureStdout(t *testing.T, ci <-chan bool) <-chan string {
	t.Helper()

	stdout := os.Stdout

	r, w, err := os.Pipe()

	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = w

	t.Cleanup(func() {
		os.Stdout = stdout
	})

	co := make(chan string)

	go func() {
		<-ci

		w.Close()

		buf := &bytes.Buffer{}

		if _, err := buf.ReadFrom(r); err != nil {
			t.Error(err)
		}

		co <- strings.TrimRight(buf.String(), "\n")
	}()

	return co
}

func Test_summary_print(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		returnCode alert.ReturnCode
		msg        string
	}

	type testCase struct {
		name     string
		args     args
		expected string
	}

	testCases := []testCase{
		{
			name: "ok",
			args: args{
				returnCode: alert.OK,
				msg:        "ok",
			},
			expected: "CLOUDWATCH OK: ok",
		},
		{
			name: "warning",
			args: args{
				returnCode: alert.Warning,
				msg:        "warn",
			},
			expected: "CLOUDWATCH WARNING: warn",
		},
		{
			name: "critical",
			args: args{
				returnCode: alert.Critical,
				msg:        "crit",
			},
			expected: "CLOUDWATCH CRITICAL: crit",
		},
		{
			name: "unknown",
			args: args{
				returnCode: alert.Unknown,
				msg:        "unkn",
			},
			expected: "CLOUDWATCH UNKNOWN: unkn",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := newSummary(true, 0)

			ci := make(chan bool)

			co := captureStdout(t, ci)

			s.print(tc.args.returnCode, tc.args.msg)

			ci <- true

			output := <-co

			assert.Equal(tc.expected, output, "output message")
		})
	}
}

func Test_summary_build(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		warnRange           string
		criticalRange       string
		datapointsThreshold string
		metricName          string
		value               float64
		timestamp           time.Time
		isWarn              bool
		isCritical          bool
		outOfWarnRange      int
		outOfCriticalRange  int
	}

	type testCase struct {
		name     string
		args     args
		expected []string
	}

	testCases := []testCase{
		{
			name: "ok",
			args: args{
				warnRange:           "0:1",
				criticalRange:       "0:2",
				datapointsThreshold: "3/4",
				metricName:          "m1",
				value:               0.1,
				timestamp:           time.Date(2022, time.January, 2, 3, 4, 5, 0, time.UTC),
				isWarn:              false,
				isCritical:          false,
				outOfWarnRange:      2,
				outOfCriticalRange:  3,
			},
			expected: []string{
				"m1 = 0.1 | value=0.1;0:1;0:2;;",
				"m1 = 0.1 @ 2022-01-02 03:04:05 +0000 UTC; above thresholds [warn,crit] = 2,3; threshold = 3/4 | value=0.1;0:1;0:2;; datapoints_warn=2;3/4;;; datapoints_crit=3;;3/4;;",
			},
		},
		{
			name: "warning",
			args: args{
				warnRange:           "1:2",
				criticalRange:       "",
				datapointsThreshold: "5/6",
				metricName:          "m2",
				value:               2.1,
				timestamp:           time.Date(2022, time.February, 3, 4, 5, 6, 0, time.UTC),
				isWarn:              true,
				isCritical:          false,
				outOfWarnRange:      5,
				outOfCriticalRange:  1,
			},
			expected: []string{
				"m2 = 2.1; above thresholds = 5 | value=2.1;1:2;;; datapoints_warn=5;5/6;;;",
				"m2 = 2.1 @ 2022-02-03 04:05:06 +0000 UTC; above thresholds [warn,crit] = 5,1; threshold = 5/6 | value=2.1;1:2;;; datapoints_warn=5;5/6;;; datapoints_crit=1;;5/6;;",
			},
		},
		{
			name: "critical",
			args: args{
				warnRange:           "",
				criticalRange:       "3:4",
				datapointsThreshold: "6/7",
				metricName:          "m3",
				value:               2.9,
				timestamp:           time.Date(2022, time.March, 4, 5, 6, 7, 0, time.UTC),
				isWarn:              false,
				isCritical:          true,
				outOfWarnRange:      1,
				outOfCriticalRange:  6,
			},
			expected: []string{
				"m3 = 2.9; above thresholds = 6 | value=2.9;;3:4;; datapoints_crit=6;;6/7;;",
				"m3 = 2.9 @ 2022-03-04 05:06:07 +0000 UTC; above thresholds [warn,crit] = 1,6; threshold = 6/7 | value=2.9;;3:4;; datapoints_warn=1;6/7;;; datapoints_crit=6;;6/7;;",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(
				tc.expected[0],
				newSummary(true, 0).build(
					tc.args.warnRange,
					tc.args.criticalRange,
					tc.args.datapointsThreshold,
					tc.args.metricName,
					tc.args.value,
					tc.args.timestamp,
					tc.args.isWarn,
					tc.args.isCritical,
					tc.args.outOfWarnRange,
					tc.args.outOfCriticalRange,
				),
				"verbosity = 0",
			)

			assert.Equal(
				tc.expected[1],
				newSummary(true, 1).build(
					tc.args.warnRange,
					tc.args.criticalRange,
					tc.args.datapointsThreshold,
					tc.args.metricName,
					tc.args.value,
					tc.args.timestamp,
					tc.args.isWarn,
					tc.args.isCritical,
					tc.args.outOfWarnRange,
					tc.args.outOfCriticalRange,
				),
				"verbosity = 1",
			)
		})
	}
}
