package alert

import (
	"fmt"

	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/errors"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/log"
)

type Checker struct {
	threshold          threshold
	isWarn             bool
	isCritical         bool
	outOfWarnRange     int
	outOfCriticalRange int
}

func NewChecker(warnRange string, criticalRange string, datapointsThreshold string) (Checker, error) {
	log.V(3).Trace().
		Str("package", "alert").
		Msg("parsing thresholds")

	threshold, err := newThreshold(warnRange, criticalRange, datapointsThreshold)

	if err != nil {
		return Checker{}, err
	}

	return Checker{
		threshold:          threshold,
		isWarn:             false,
		isCritical:         false,
		outOfWarnRange:     0,
		outOfCriticalRange: 0,
	}, nil
}

func (c *Checker) CheckStatus(values []float64) (ReturnCode, error) {
	log.V(3).Trace().
		Str("package", "alert").
		Msg("checking if metrics are above thresholds")

	if len(values) < c.threshold.evaluationPeriods {
		return Unknown, errors.NewArgumentErrorWithMessage(
			fmt.Sprintf("insufficient number of metrics to evaluate: got %d datapoints", len(values)),
			"datapoints",
			fmt.Sprintf("%d/%d", c.threshold.datapointsToAlarm, c.threshold.evaluationPeriods),
		)
	}

	warnCounter := newCounter(c.threshold.warn, c.threshold.datapointsToAlarm)
	criticalCounter := newCounter(c.threshold.critical, c.threshold.datapointsToAlarm)

	for i := 0; i < c.threshold.evaluationPeriods; i++ {
		warnCounter.examine(values[i])
		criticalCounter.examine(values[i])
	}

	c.outOfWarnRange = warnCounter.count
	c.outOfCriticalRange = criticalCounter.count

	if criticalCounter.over() {
		log.V(3).Trace().
			Str("package", "alert").
			Str("status", "critical").
			Int("out_of_warn_range", c.outOfWarnRange).
			Int("out_of_critical_range", c.outOfCriticalRange).
			Int("datapoints_to_alarm", c.threshold.datapointsToAlarm).
			Msg("service status is unhealthy")

		c.isCritical = true

		return Critical, nil
	}

	if warnCounter.over() {
		log.V(3).Trace().
			Str("package", "alert").
			Str("status", "warn").
			Int("out_of_warn_range", c.outOfWarnRange).
			Int("out_of_critical_range", c.outOfCriticalRange).
			Int("datapoints_to_alarm", c.threshold.datapointsToAlarm).
			Msg("service status is unhealthy")

		c.isWarn = true

		return Warning, nil
	}

	log.V(3).Trace().
		Str("package", "alert").
		Str("status", "ok").
		Int("out_of_warn_range", c.outOfWarnRange).
		Int("out_of_critical_range", c.outOfCriticalRange).
		Int("datapoints_to_alarm", c.threshold.datapointsToAlarm).
		Msg("service status is healthy")

	return OK, nil
}

func (c Checker) Result() (isWarn bool, isCritical bool, outOfWarnRange int, outOfCriticalRange int) {
	isWarn = c.isWarn
	isCritical = c.isCritical

	outOfWarnRange = c.outOfWarnRange
	outOfCriticalRange = c.outOfCriticalRange

	return
}
