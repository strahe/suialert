package storage

import (
	"context"
	"fmt"
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	"github.com/strahe/suialert/model/schema"
)

// initDatabaseSchema initializes the version tables for tracking schema version installed in the database
func initDatabaseSchema(ctx context.Context, db *pg.DB, cfg schema.Config) error {
	if cfg.SchemaName != "public" {
		_, err := db.Exec(`CREATE SCHEMA IF NOT EXISTS ?`, pg.SafeQuery(cfg.SchemaName))
		if err != nil {
			return fmt.Errorf("ensure schema exists :%w", err)
		}
	}

	// Ensure the pg migrations table exists
	migTableName := cfg.SchemaName + ".pg_migrations"
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ? (
			id serial,
			version bigint,
			created_at timestamptz
		)
	`, pg.SafeQuery(migTableName))
	if err != nil {
		return fmt.Errorf("ensure visor_version exists :%w", err)
	}

	return createSchema(db)
}

func collectionForVersion(cfg schema.Config) (*migrations.Collection, error) {
	return schema.GetPatches(cfg)
}

func getDatabaseSchemaVersion(ctx context.Context, db *pg.DB, cfg schema.Config) (int, bool, error) {
	migExists, err := tableExists(ctx, db, cfg.SchemaName, "pg_migrations")
	if err != nil {
		return 0, false, fmt.Errorf("checking if pg_migrations exists:%w", err)
	}

	if !migExists {
		// Uninitialized database
		return 0, false, nil
	}

	coll, err := collectionForVersion(cfg)
	if err != nil {
		return 0, false, err
	}

	migration, err := coll.Version(db)
	if err != nil {
		return 0, false, fmt.Errorf("unable to determine schema version: %w", err)
	}

	if migration == 0 {
		return 0, false, nil
	}
	return int(migration), true, nil
}

func validateDatabaseSchemaVersion(ctx context.Context, db *pg.DB, cfg schema.Config) error {
	// Check if the version of the schema is compatible
	_, initialized, err := getDatabaseSchemaVersion(ctx, db, cfg)
	if err != nil {
		return fmt.Errorf("get schema version: %w", err)
	}

	if !initialized {
		return fmt.Errorf("schema not installed in database")
	}

	return nil
}
