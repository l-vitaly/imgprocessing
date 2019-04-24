package log

import (
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// LoggerWrapDefault wrap for default options.
func LoggerWrapDefault(logger kitlog.Logger, service, levelEnvName string) kitlog.Logger {
	logger = kitlog.With(logger, "@service", service)
	logger = kitlog.With(logger, "@timestamp", kitlog.DefaultTimestampUTC)
	logger = kitlog.With(logger, "caller", kitlog.Caller(5))
	return level.NewFilter(logger, getLevelOption(levelEnvName))
}

func getLevelOption(key string) level.Option {
	switch os.Getenv(key) {
	case "error":
		return level.AllowError()
	case "info":
		return level.AllowInfo()
	case "warn":
		return level.AllowWarn()
	default:
		return level.AllowDebug()
	}
}
