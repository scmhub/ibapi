package ibapi

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

const DEFAULT_LEVEL = zerolog.InfoLevel

func init() {
	zerolog.SetGlobalLevel(DEFAULT_LEVEL)
	log = zerolog.New(os.Stderr).With().Timestamp().Logger()
}

// Logger returns the logger.
func Logger() *zerolog.Logger {
	return &log
}

// SetLogLevel sets the loggging level.
func SetLogLevel(logLevel int) {
	zerolog.SetGlobalLevel(zerolog.Level(int8(logLevel)))
}

// SetConsoleWriter will send pretty log to the console.
func SetConsoleWriter() {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("| IB | %s", i)
	}
	log = log.Output(output)
}
