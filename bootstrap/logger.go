package bootstrap

import (
	"os"

	"github.com/rs/zerolog"
)

func InitLogger(version string, debug bool) zerolog.Logger {
	host, _ := os.Hostname()
	logger := zerolog.
		New(os.Stderr).
		With().
		Caller().
		Timestamp().
		Str("host", host).
		Str("version", version).
		Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	return logger
}
