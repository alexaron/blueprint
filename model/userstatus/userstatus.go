// Package userstatus provides access to the user_status table in the MySQL database.
package userstatus

import (
	"github.com/gocraft/dbr"
)

var (
	// table is the table name.
	table = "user_status"
)

// Item defines the model
type Item struct {
	ID        uint8        `db:"id"`
	Status    string       `db:"status"`
	CreatedAt dbr.NullTime `db:"created_at"`
	UpdatedAt dbr.NullTime `db:"updated_at"`
	DeletedAt dbr.NullTime `db:"deleted_at"`
}
