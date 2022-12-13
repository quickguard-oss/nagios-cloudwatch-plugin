package main

import (
	"fmt"
	"time"

	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/alert"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/log"
)

type summary struct {
	classicOutput bool
	isVerbose     bool
}

const pluginName string = "CLOUDWATCH"

func newSummary(classicOutput bool, verbosity int) summary {
	log.V(3).Trace().
		Str("package", "main").
		Bool("classic_output", classicOutput).
		Int("verbosity", verbosity).
		Msg("set output options")

	return summary{
		classicOutput: classicOutput,
		isVerbose:     1 <= verbosity,
	}
}

func (o summary) print(returnCode alert.ReturnCode, msg string) {
	if o.classicOutput {
		fmt.Printf("%s %s: %s\n", pluginName, returnCode.String(), msg)
	} else {
		log.V(0).Info().
			Str("service", pluginName).
			Str("status", returnCode.String()).
			Msg(msg)
	}
}

func (o summary) build(
	warnRange string, criticalRange string, datapointsThreshold string,
	metricName string, value float64, timestamp time.Time,
	isWarn bool, isCritical bool, outOfWarnRange int, outOfCriticalRange int,
) string {
	if o.isVerbose {
		return fmt.Sprintf(
			"%s = %g @ %s; above thresholds [warn,crit] = %d,%d; threshold = %s | value=%g;%s;%s;; datapoints_warn=%d;%s;;; datapoints_crit=%d;;%s;;",
			metricName, value, timestamp, outOfWarnRange, outOfCriticalRange, datapointsThreshold,
			value, warnRange, criticalRange, outOfWarnRange, datapointsThreshold, outOfCriticalRange, datapointsThreshold,
		)
	} else {
		if isCritical {
			return fmt.Sprintf(
				"%s = %g; above thresholds = %d | value=%g;%s;%s;; datapoints_crit=%d;;%s;;",
				metricName, value, outOfCriticalRange,
				value, warnRange, criticalRange, outOfCriticalRange, datapointsThreshold,
			)
		} else if isWarn {
			return fmt.Sprintf(
				"%s = %g; above thresholds = %d | value=%g;%s;%s;; datapoints_warn=%d;%s;;;",
				metricName, value, outOfWarnRange,
				value, warnRange, criticalRange, outOfWarnRange, datapointsThreshold,
			)
		} else {
			return fmt.Sprintf(
				"%s = %g | value=%g;%s;%s;;",
				metricName, value,
				value, warnRange, criticalRange,
			)
		}
	}
}
