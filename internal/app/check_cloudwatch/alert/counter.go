package alert

import (
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/log"
)

type counter struct {
	count             int
	thresholdRange    thresholdRange
	datapointsToAlarm int
}

func newCounter(t thresholdRange, n int) counter {
	return counter{
		count:             0,
		thresholdRange:    t,
		datapointsToAlarm: n,
	}
}

func (c *counter) examine(value float64) {
	if !c.thresholdRange.enable {
		return
	}

	if c.outOfRange(value) {
		log.V(3).Trace().
			Str("package", "alert").
			Bool("above_threshold", true).
			Float64("value", value).
			Float64("range_start", c.thresholdRange.start).
			Float64("range_end", c.thresholdRange.end).
			Bool("alert_if_inside_range", c.thresholdRange.inverse).
			Msg("the value is above threshold")

		c.increment()
	} else {
		log.V(3).Trace().
			Str("package", "alert").
			Bool("above_threshold", false).
			Float64("value", value).
			Float64("range_start", c.thresholdRange.start).
			Float64("range_end", c.thresholdRange.end).
			Bool("alert_if_inside_range", c.thresholdRange.inverse).
			Msg("the value is below threshold")
	}
}

func (c counter) outOfRange(value float64) bool {
	isOutside := (value < c.thresholdRange.start) || (c.thresholdRange.end < value)

	if c.thresholdRange.inverse {
		return !isOutside
	} else {
		return isOutside
	}
}

func (c *counter) increment() {
	c.count++
}

func (c counter) over() bool {
	return c.datapointsToAlarm <= c.count
}
