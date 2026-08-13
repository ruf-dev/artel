// SPDX-License-Identifier: AGPL-3.0-only
package main

import (
	"github.com/rs/zerolog/log"

	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/config"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	err = postgres.Migrate(cfg.DataSources.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to apply migrations")
	}

	log.Info().Msg("migrations applied successfully")
}
