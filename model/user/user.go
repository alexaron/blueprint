// Package user provides access to the user table in the MySQL database.
package user

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"

	"github.com/go-sql-driver/mysql"
)

var (
	// table is the table name.
	table = "user"
)

// Item defines the model.
type Item struct {
	ID               uint32         `db:"id"`
	FirstName        string         `db:"first_name"`
	LastName         string         `db:"last_name"`
	Email            string         `db:"email"`
	Password         string         `db:"password"`
	VerificationCode string         `db:"verification_code"`
	Verified         bool           `db:"verified"`
	StatusID         uint8          `db:"status_id"`
	CreatedAt        mysql.NullTime `db:"created_at"`
	UpdatedAt        mysql.NullTime `db:"updated_at"`
	DeletedAt        mysql.NullTime `db:"deleted_at"`
}

// Connection is an interface for making queries.
type Connection interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

// ByEmail gets user information from email.
func ByEmail(db Connection, email string) (Item, bool, error) {
	result := Item{}
	err := db.Get(&result, fmt.Sprintf(`
		SELECT id, password, status_id, first_name, verification_code, verified
		FROM %v
		WHERE email = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		email)
	return result, err == sql.ErrNoRows, err
}

// ByEmail gets user information from verification_code.
func ByCode(db Connection, code string) (Item, bool, error) {
	result := Item{}
	err := db.Get(&result, fmt.Sprintf(`
		SELECT id, password, status_id, first_name, verification_code, verified
		FROM %v
		WHERE verification_code = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		code)
	return result, err == sql.ErrNoRows, err
}

// Create creates user.
func Create(db Connection, firstName, lastName, email, password string) (sql.Result, error) {
	code := pseudoSha2()
	result, err := db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(first_name, last_name, email, password, verification_code)
		VALUES
		(?,?,?,?,?)
		`, table),
		firstName, lastName, email, password, code)
	return result, err
}

// Verify marks a user as verified.
func Verify(db Connection, id uint32) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		UPDATE %v
		SET verification_code='', verified=1
		WHERE id=?
		`, table),
		id)
	return result, err
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
