package db

import (
	"database/sql"

	"github.com/blue-jay/blueprint/lib/env"
	"github.com/gocraft/dbr"
	"github.com/gocraft/dbr/dialect"
	_ "github.com/jackc/pgx/stdlib"

	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

// Connect connects to the database specified in the config.
func Connect(conf env.PostgreSQL) (*dbr.Connection, error) {
	conn, err := sql.Open("pgx", conf.DSN())
	if err != nil {
		return nil, err
	}
	return &dbr.Connection{
		DB:            conn,
		Dialect:       dialect.PostgreSQL,
		EventReceiver: &dbr.NullEventReceiver{},
	}, nil
}

// MigrateUp runs migrations against the provided config. The database must already exist.
func MigrateUp(dbSess *dbr.Session, dbName string) error {
	driver, err := postgres.WithInstance(dbSess.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migration/postgres",
		dbName, driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != migrate.ErrNoChange {
		return err
	}
	return nil
}
