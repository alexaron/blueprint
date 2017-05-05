// Package user provides access to the user table in the MySQL database.
package user

import (
	"crypto/rand"
	"database/sql"
	"math/big"

	"github.com/gocraft/dbr"
)

var (
	// table is the table name.
	table = "user"
)

// Item defines the model.
type Item struct {
	ID               uint32       `db:"id"`
	FirstName        string       `db:"first_name"`
	LastName         string       `db:"last_name"`
	Email            string       `db:"email"`
	Password         string       `db:"password"`
	VerificationCode string       `db:"verification_code"`
	Verified         bool         `db:"verified"`
	StatusID         uint8        `db:"status_id"`
	CreatedAt        dbr.NullTime `db:"created_at"`
	UpdatedAt        dbr.NullTime `db:"updated_at"`
	DeletedAt        dbr.NullTime `db:"deleted_at"`
}

// ByEmail gets user information from email.
func ByEmail(db *dbr.Session, email string) (Item, bool, error) {
	result := Item{}
	err := db.
		Select("*").
		From(dbr.I(table)). // user is a reserved word in PostgreSQL, needs quoting.
		Where("email = ? AND deleted_at IS NULL", email).
		Limit(1).
		LoadStruct(&result)
	return result, err == dbr.ErrNotFound, err
}

// ByEmail gets user information from verification_code.
func ByCode(db *dbr.Session, code string) (Item, bool, error) {
	result := Item{}
	err := db.
		Select("*").
		From(dbr.I(table)).
		Where("verification_code = ? AND deleted_at IS NULL", code).
		Limit(1).
		LoadStruct(&result)
	return result, err == dbr.ErrNotFound, err
}

// Create creates user.
func Create(db *dbr.Session, firstName, lastName, email, password string) (sql.Result, error) {
	// Using raw SQL statement because LastInsertID is used in the outside
	// and PostgreSQL only supports it when the RETURNING statement is provided.
	return db.InsertBySql(`
		INSERT INTO "`+table+`"
			(first_name, last_name, email, password, verification_code)
			VALUES (?, ?, ?, ?, ?)
			RETURNING id
	`, firstName, lastName, email, password, pseudoSha2()).Exec()
}

// Verify marks a user as verified.
func Verify(db *dbr.Session, id uint32) (sql.Result, error) {
	return db.
		Update(table).
		Set("verification_code", "").
		Set("verified", true).
		Where("id = ?", id).
		Exec()
}

// pseudoSha2 outputs 64-byte string that looks like the real sha2.
func pseudoSha2() string {
	hexNums := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
	maxRand := big.NewInt(16)

	l := 64
	result := make([]byte, l)
	for i := 0; i < l; i++ {
		// The error is checked in tests.
		t, _ := rand.Int(rand.Reader, maxRand)
		result[i] = hexNums[t.Int64()]
	}
	return string(result)
}
