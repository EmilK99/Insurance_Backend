package store

import (
	"context"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"

	"github.com/jackc/tern/migrate"
)

type Store struct {
	Conn *pgx.Conn
}

func NewStore(conn *pgx.Conn) *Store {
	return &Store{Conn: conn}
}

func MigrateDatabase(ctx context.Context, conn *pgx.Conn) {
	migrator, err := migrate.NewMigrator(ctx, conn, "schema_migration")
	if err != nil {
		log.Fatalf("Unable to create a migrator: %v", err)
	}

	err = migrator.LoadMigrations("./app/store/migrations")
	if err != nil {
		log.Fatalf("Unable to load migrations: %v", err)
	}

	err = migrator.Migrate(ctx)

	if err != nil {
		log.Fatalf("Unable to migrate: %v", err)
	}

	ver, err := migrator.GetCurrentVersion(ctx)
	if err != nil {
		log.Fatalf("Unable to get current schema version: %v", err)
	}

	log.Infof("Migration done. Current schema version: %v", ver)
}
