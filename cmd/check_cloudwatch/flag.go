package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/errors"
	"github.com/spf13/pflag"
)

type flags struct {
	warnRange           *string
	criticalRange       *string
	datapointsThreshold *string
	queries             *string
	duration            *int
	timeout             *int
	classicOutput       *bool
	verbosity           *int
	showVersion         *bool
	showHelp            *bool
}

func newFlags() flags {
	return flags{}
}

func parseFlags() (flags, error) {
	setupParser()

	f := newFlags()

	f.defineFlags()

	if err := pflag.CommandLine.Parse(os.Args[1:]); err != nil {
		return f, errors.NewArgumentErrorWithError(err, "arguments", strings.Join(os.Args[1:], " "))
	}

	if *f.showVersion {
		fmt.Println(version)

		os.Exit(0)
	}

	if *f.showHelp {
		pflag.Usage()

		os.Exit(0)
	}

	if *f.queries == "" {
		return f, errors.NewArgumentErrorWithMessage("queries must be an array of MetricDataQuery objects", "queries", "")
	}

	if *f.duration <= 0 {
		return f, errors.NewArgumentErrorWithMessage("time duration must be a positive number", "duration", strconv.Itoa(*f.duration))
	}

	if *f.timeout <= 0 {
		return f, errors.NewArgumentErrorWithMessage("timeout must be a positive number", "timeout", strconv.Itoa(*f.timeout))
	}

	return f, nil
}

func setupParser() {
	pflag.CommandLine.Init(os.Args[0], pflag.ContinueOnError)

	pflag.CommandLine.SetOutput(os.Stdout)

	pflag.CommandLine.SortFlags = false

	pflag.Usage = func() {
		header := fmt.Sprintf("check_cloudwatch (v%s)\n", version)

		usage := `
This plugin checks AWS CloudWatch metrics using GetMetricData API.

Usage:
  check_cloudwatch -q <queries> -w <range> -c <range> -p <datapoints>
                   [-d <duration>] [-t <timeout>] [-C] [-v]

Options:
`

		fmt.Printf(header + usage)

		pflag.PrintDefaults()
	}
}

func (f *flags) defineFlags() {
	f.queries = pflag.StringP(
		"queries", "q",
		"",
		"Array of MetricDataQuery objects in JSON format. See AWS GetMetricData API reference.",
	)

	f.warnRange = pflag.StringP(
		"warning", "w",
		"",
		"Metric range to result in warning status.",
	)

	f.criticalRange = pflag.StringP(
		"critical", "c",
		"",
		"Metric range to result in critical status.",
	)

	f.datapointsThreshold = pflag.StringP(
		"datapoints", "p",
		"1/1",
		"'n/m'; resulting to anomaly status if n of m datapoints are out of warning or critical range.",
	)

	f.duration = pflag.IntP(
		"duration", "d",
		60,
		"Duration in minutes for which to retrieve metrics.",
	)

	f.timeout = pflag.IntP(
		"timeout", "t",
		10,
		"Seconds before plugin times out.",
	)

	f.classicOutput = pflag.BoolP(
		"classic-output", "C",
		false,
		"Print status message in classic format.",
	)

	f.verbosity = pflag.CountP(
		"verbose", "v",
		"Extra information. Up to 3 verbosity levels.",
	)

	f.showVersion = pflag.BoolP(
		"version", "V",
		false,
		"Print version information.",
	)

	f.showHelp = pflag.BoolP(
		"help", "h",
		false,
		"Print detailed help screen.",
	)
}
