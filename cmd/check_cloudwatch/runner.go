package main

import (
	"time"

	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/alert"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/cloudwatch"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/log"
)

func run() alert.ReturnCode {
	flags, err := parseFlags()

	log.SetVerbosity(*flags.verbosity)

	summary := newSummary(*flags.classicOutput, *flags.verbosity)

	if err != nil {
		summary.print(alert.Unknown, err.Error())

		return alert.Unknown
	}

	checker, err := alert.NewChecker(*flags.warnRange, *flags.criticalRange, *flags.datapointsThreshold)

	if err != nil {
		summary.print(alert.Unknown, err.Error())

		return alert.Unknown
	}

	client, err := cloudwatch.New(*flags.duration, *flags.queries, *flags.timeout)

	if err != nil {
		summary.print(alert.Unknown, err.Error())

		return alert.Unknown
	}

	values, err := client.GetMetricValues(time.Now())

	if err != nil {
		summary.print(alert.Unknown, err.Error())

		return alert.Unknown
	}

	returnCode, err := checker.CheckStatus(values)

	if err != nil {
		summary.print(returnCode, err.Error())
	} else {
		a1, a2, a3 := client.LatestValue()
		b1, b2, b3, b4 := checker.Result()

		summary.print(
			returnCode,
			summary.build(
				*flags.warnRange, *flags.criticalRange, *flags.datapointsThreshold,
				a1, a2, a3, b1, b2, b3, b4,
			),
		)
	}

	return returnCode
}
