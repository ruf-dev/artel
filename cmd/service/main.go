// SPDX-License-Identifier: AGPL-3.0-only
package main

import (
	"github.com/rs/zerolog/log"
	"github.com/ruf-dev/artel/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create application")
	}

	err = a.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create application")
	}
}
