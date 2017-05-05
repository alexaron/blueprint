package note_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/blue-jay/blueprint/model/note"
	"github.com/blue-jay/blueprint/model/user"

	"github.com/blue-jay/blueprint/lib/db"
	"github.com/blue-jay/blueprint/lib/env"
	"github.com/gocraft/dbr"
)

var dbTestSess *dbr.Session

const (
	TestDatabase = "_test_blueprint"
)

// TestMain runs setup, tests, and then teardown.
func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

// setup handles any start up tasks.
func setup() {
	info, err := env.LoadConfig("../../env.json")
	if err != nil {
		panic("Error loading config: " + err.Error())
	}

	// Connecting to the existing database and creating a test one.
	dbConn, err := db.Connect(info.PostgreSQL)
	if err != nil {
		panic("Error connecting to the existing database: " + err.Error())
	}

	// First - drop an existing TestDatabase and then create a new one.
	// This way you'll have results after running tests at hand.
	if _, err := dbConn.Exec("DROP DATABASE IF EXISTS " + TestDatabase); err != nil {
		panic("Error dropping test database: " + err.Error())
	}
	if _, err := dbConn.Exec("CREATE DATABASE " + TestDatabase); err != nil {
		panic("Error creating test database: " + err.Error())
	}

	// Connecting to the newly created test database.
	info.PostgreSQL.Database = TestDatabase
	testDbConn, err := db.Connect(info.PostgreSQL)
	if err != nil {
		panic("Error connecting to TestDatabase: " + err.Error())
	}
	dbTestSess = testDbConn.NewSession(nil)

	// And running migrations.
	if err := db.MigrateUp(dbTestSess, info.PostgreSQL.Database); err != nil {
		panic("Error running migrations: " + err.Error())
	}
}

// TestComplete
func TestComplete(t *testing.T) {
	data := "Test data."
	dataNew := "New test data."

	result, err := user.Create(dbTestSess, "John", "Doe", "jdoe@domain.com", "p@$$W0rD")
	if err != nil {
		t.Fatal("could not create user:", err)
	}

	uID, err := result.LastInsertId()
	if err != nil {
		t.Fatal("could not convert user ID:", err)
	}

	// Convert ID to string
	userID := fmt.Sprintf("%v", uID)

	// Create a record
	result, err = note.Create(dbTestSess, data, userID)
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
	record, _, err := note.ByID(dbTestSess, lastID, userID)
	if err != nil {
		t.Error("could not retrieve record:", err)
	} else if record.Name != data {
		t.Errorf("retrieved wrong record: got '%v' want '%v'", record.Name, data)
	}

	// Update a record
	result, err = note.Update(dbTestSess, dataNew, lastID, userID)
	if err != nil {
		t.Error("could not update record:", err)
	}

	// Select a record
	record, _, err = note.ByID(dbTestSess, lastID, userID)
	if err != nil {
		t.Error("could not retrieve record:", err)
	} else if record.Name != dataNew {
		t.Errorf("retieved wrong record: got '%v' want '%v'", record.Name, dataNew)
	}

	// Delete a record by ID
	result, err = note.DeleteSoft(dbTestSess, lastID, userID)
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
