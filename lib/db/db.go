package db

import (
	"github.com/gocraft/dbr"
	"github.com/gocraft/dbr/dialect"

	"github.com/blue-jay/core/storage/driver/mysql"
)

// Connect connects to the database specified in the config.
// FIXME: Refactor this to use only dbr without sqlx. sqlx is only left for testing purposes.
func Connect(conf mysql.Info, specificDatabase bool) (*dbr.Connection, error) {
	sqlxConn, err := conf.Connect(specificDatabase)
	if err != nil {
		return nil, err
	}
	return &dbr.Connection{
		DB:            sqlxConn.DB,
		Dialect:       dialect.MySQL,
		EventReceiver: &dbr.NullEventReceiver{},
	}, nil
}
