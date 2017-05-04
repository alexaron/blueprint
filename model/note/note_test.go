package note_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/blue-jay/blueprint/model/note"
	"github.com/blue-jay/blueprint/model/user"
	"github.com/blue-jay/core/storage/migration/mysql"

	"github.com/blue-jay/blueprint/lib/db"
	"github.com/gocraft/dbr"
	"github.com/jmoiron/sqlx"
)

var (
	testDb *dbr.Session
)

// TestMain runs setup, tests, and then teardown.
func TestMain(m *testing.M) {
	setup()
	returnCode := m.Run()
	teardown()
	os.Exit(returnCode)
}

// setup handles any start up tasks.
func setup() {
	_, conf := mysql.SetUp("../../env.json.example", "database_test")

	// Connect to the database
	d, _ := db.Connect(conf.Info, true)
	testDb = d.NewSession(nil)
}

// teardown handles any clean up tasks.
func teardown() {
	mysql.TearDown(sqlx.NewDb(testDb.DB, "mysql"), "database_test")
}

// TestComplete
func TestComplete(t *testing.T) {
	data := "Test data."
	dataNew := "New test data."

	result, err := user.Create(testDb, "John", "Doe", "jdoe@domain.com", "p@$$W0rD")
	if err != nil {
		t.Error("could not create user:", err)
	}

	uID, err := result.LastInsertId()
	if err != nil {
		t.Error("could not convert user ID:", err)
	}

	// Convert ID to string
	userID := fmt.Sprintf("%v", uID)

	// Create a record
	result, err = note.Create(testDb, data, userID)
	if err != nil {
		t.Error("could not create record:", err)
	}

	// Get the last ID
	ID, err := result.LastInsertId()
	if err != nil {
		t.Error("could not convert ID:", err)
	}

	// Convert ID to string
	lastID := fmt.Sprintf("%v", ID)

	// Select a record
	record, _, err := note.ByID(testDb, lastID, userID)
	if err != nil {
		t.Error("could not retrieve record:", err)
	} else if record.Name != data {
		t.Errorf("retrieved wrong record: got '%v' want '%v'", record.Name, data)
	}

	// Update a record
	result, err = note.Update(testDb, dataNew, lastID, userID)
	if err != nil {
		t.Error("could not update record:", err)
	}

	// Select a record
	record, _, err = note.ByID(testDb, lastID, userID)
	if err != nil {
		t.Error("could not retrieve record:", err)
	} else if record.Name != dataNew {
		t.Errorf("retieved wrong record: got '%v' want '%v'", record.Name, dataNew)
	}

	// Delete a record by ID
	result, err = note.DeleteSoft(testDb, lastID, userID)
	if err != nil {
		t.Error("could not delete record:", err)
	}

	// Count the number of deleted rows
	rows, err := result.RowsAffected()
	if err != nil {
		t.Error("could not count affected rows:", err)
	} else if rows != 1 {
		t.Error("incorrect number of affected rows:", rows)
	}
}
