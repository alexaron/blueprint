package db

import (
	"github.com/gocraft/dbr"

	"database/sql"
	"github.com/blue-jay/blueprint/lib/env"
	"github.com/gocraft/dbr/dialect"
	_ "github.com/jackc/pgx/stdlib"
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
