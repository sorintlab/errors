package main

import (
	"os"
	"time"

	"github.com/sorintlab/errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// setup zerolog console writer using our custom FormatErrFieldValue
	cw := zerolog.ConsoleWriter{
		Out:                 os.Stderr,
		TimeFormat:          time.RFC3339Nano,
		FormatErrFieldValue: errors.FormatErrFieldValue,
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano

	// setup a custom zerolog ErrorMarshalFunc to save detailed error data
	zerolog.ErrorMarshalFunc = errors.ErrorMarshalFunc

	log.Logger = log.With().Caller().Logger().Level(zerolog.InfoLevel).Output(cw)
}

func main() {
	err := errors.Errorf("initial error")
	err = errors.Wrapf(err, "wrapped error")
	err = errors.WithStack(err)
	log.Err(err).Msg("there was an error!")
}
