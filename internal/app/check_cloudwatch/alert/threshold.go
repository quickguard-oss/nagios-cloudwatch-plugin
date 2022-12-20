package alert

import (
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/errors"
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/log"
)

type threshold struct {
	evaluationPeriods int
	datapointsToAlarm int
	warn              thresholdRange
	critical          thresholdRange
}

type thresholdRange struct {
	enable  bool
	start   float64
	end     float64
	inverse bool
}

func newThreshold(warnRange string, criticalRange string, datapointsThreshold string) (threshold, error) {
	log.V(3).Trace().
		Str("package", "alert").
		Msg("parsing warn threshold")

	warn, err := newThresholdRange(warnRange)

	if err != nil {
		return threshold{}, err
	}

	log.V(3).Trace().
		Str("package", "alert").
		Msg("parsing critical threshold")

	critical, err := newThresholdRange(criticalRange)

	if err != nil {
		return threshold{}, err
	}

	log.V(3).Trace().
		Str("package", "alert").
		Msg("parsing datapoints threshold")

	datapointsToAlarm, evaluationPeriods, err := parseDatapointsThreshold(datapointsThreshold)

	if err != nil {
		return threshold{}, err
	}

	return threshold{
		evaluationPeriods: evaluationPeriods,
		datapointsToAlarm: datapointsToAlarm,
		warn:              warn,
		critical:          critical,
	}, nil
}

func newThresholdRange(s string) (thresholdRange, error) {
	var start, end float64
	var startStr, endStr string

	if s == "" {
		log.V(3).Trace().
			Str("package", "alert").
			Float64("range_start", math.Inf(-1)).
			Float64("range_end", math.Inf(1)).
			Bool("alert_if_inside_range", false).
			Msg("range is not specified; use default value")

		return thresholdRange{}, nil
	}

	num := `-?\d+(?:\.\d+)?`

	result := regexp.MustCompile(
		fmt.Sprintf(`\A(@)?(?:(%s|~):)?(%s)?\z`, num, num),
	).FindStringSubmatch(s)

	result = append(result, make([]string, 4-len(result))...)

	switch {
	case result[2] != "" && result[3] != "":
		startStr = result[2]
		endStr = result[3]
	case result[2] != "" && result[3] == "":
		startStr = result[2]
		endStr = "~"
	case result[2] == "" && result[3] != "":
		startStr = "0"
		endStr = result[3]
	case result[2] == "" && result[3] == "":
		return thresholdRange{}, errors.NewArgumentErrorWithMessage("warn/critical range is specified in the format 'START:END'", "warn/critical range", s)
	}

	if startStr == "~" {
		start = math.Inf(-1)
	} else {
		start, _ = strconv.ParseFloat(startStr, 64)
	}

	if endStr == "~" {
		end = math.Inf(1)
	} else {
		end, _ = strconv.ParseFloat(endStr, 64)
	}

	if end < start {
		return thresholdRange{}, errors.NewArgumentErrorWithMessage("end value must be greater than start value", "warn/critical range", s)
	}

	inverse := result[1] == "@"

	log.V(3).Trace().
		Str("package", "alert").
		Float64("range_start", start).
		Float64("range_end", end).
		Bool("alert_if_inside_range", inverse).
		Send()

	return thresholdRange{
		enable:  true,
		start:   start,
		end:     end,
		inverse: inverse,
	}, nil
}

func parseDatapointsThreshold(s string) (int, int, error) {
	result := regexp.MustCompile(`\A(\d+)/(\d+)\z`).FindStringSubmatch(s)

	result = append(result, make([]string, 3-len(result))...)

	datapointsToAlarmStr := result[1]
	evaluationPeriodsStr := result[2]

	if datapointsToAlarmStr == "" || evaluationPeriodsStr == "" {
		return 0, 0, errors.NewArgumentErrorWithMessage("evaluation periods (=m) and datapoints to alarm (=n) are specified in the format 'n/m'", "datapoints", s)
	}

	datapointsToAlarm, _ := strconv.Atoi(datapointsToAlarmStr)
	evaluationPeriods, _ := strconv.Atoi(evaluationPeriodsStr)

	if evaluationPeriods < datapointsToAlarm {
		return 0, 0, errors.NewArgumentErrorWithMessage("evaluation periods must be greater than datapoints to alarm", "datapoints", s)
	}

	log.V(3).Trace().
		Str("package", "alert").
		Int("evaluation_periods", evaluationPeriods).
		Int("datapoints_to_alarm", datapointsToAlarm).
		Send()

	return datapointsToAlarm, evaluationPeriods, nil
}
