package postgres

import (
	"github.com/betas-in/logger"
	"github.com/golang-migrate/migrate/v4"

	// Import driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// Import driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// Import driver
	_ "github.com/lib/pq"
)

// Migrate ...
func Migrate(conf *Config, logger *logger.Logger, direction string, version int) error {
	logger.Info("postgres.migrate").Msgf("Attempting to run %s migration on %s from %+v", direction, ConnectionURL(conf), conf.MigrationPath)
	m, err := migrate.New(conf.MigrationPath, ConnectionURL(conf))
	if err != nil {
		logger.Fatal("postgres.migrate").Err(err).Msg("Could not run the migration")
	}

	switch direction {
	case "up":
		// MigrateUp runs database migration
		err = m.Up()
	case "down":
		// MigrateDown rolls back the latest migration
		err = m.Down()
	case "force":
		// MigrateForce force a specific version
		err = m.Force(version)
	}
	if err == migrate.ErrNoChange || err == nil {
		return nil
	}

	logger.Fatal("postgres.migrate").Err(err).Msgf("Failed to run the %s migration", direction)
	return err
}
