package utils

import (
	"io"

	"github.com/rs/zerolog/log"
)

func CloseWithLog(c io.Closer, description string) {
	err := c.Close()
	if err != nil {
		log.Error().
			Err(err).
			Str("resource", description).
			Msg("error closing resource")
	}
}

type errorFunc func() error

func CallWithLog(f errorFunc, errMsg string) {
	err := f()
	if err != nil {
		log.Error().
			Err(err).
			Str("resource", errMsg).
			Msg("error calling resource")
	}
}
