package log

import (
	"github.com/quickguard-oss/nagios-cloudwatch-plugin/internal/app/check_cloudwatch/container"
	"github.com/rs/zerolog"
)

const maxVerbosity int = 3

type loggerSet [maxVerbosity + 1]zerolog.Logger

type loggerStatuses [maxVerbosity + 1]struct {
	initialized  bool
	defaultLevel zerolog.Level
}

var loggers loggerSet

var status = loggerStatuses{
	{defaultLevel: zerolog.TraceLevel},
	{defaultLevel: zerolog.Disabled},
	{defaultLevel: zerolog.Disabled},
	{defaultLevel: zerolog.Disabled},
}

func Reset() {
	loggers = loggerSet{}

	for i := 0; i <= maxVerbosity; i++ {
		status[i].initialized = false
	}
}

func SetVerbosity(verbosity int) {
	for i := 0; i <= maxVerbosity; i++ {
		if i <= verbosity {
			loggers[i] = getLogger(i).Level(zerolog.TraceLevel)
		} else {
			loggers[i] = getLogger(i).Level(zerolog.Disabled)
		}
	}
}

func V(verbosity int) *zerolog.Logger {
	l := getLogger(verbosity)

	return &l
}

func getLogger(verbosity int) zerolog.Logger {
	s := &status[verbosity]

	if !s.initialized {
		loggers[verbosity] = newLogger(s.defaultLevel)

		s.initialized = true
	}

	return loggers[verbosity]
}

func newLogger(logLevel zerolog.Level) zerolog.Logger {
	return zerolog.New(
		container.GetLoggerIO(),
	).Level(logLevel).With().Timestamp().Logger()
}
