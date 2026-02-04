package database

import (
	"path/filepath"
	"testing"
)

type testUser struct {
	id             int
	email          string
	password       string
	hashedPassword string
}

var testUsers = map[string]*testUser{
	"alice": {email: "alice@github.com/glenntam/todoken", password: "testPass123!", hashedPassword: "$2a$04$mi5gstbTPDRpEawTIitij.rdzLFM.U8.x4U5LLzK8xVFXKXf2ng2u"},
	"bob":   {email: "bob@github.com/glenntam/todoken", password: "mySecure456#", hashedPassword: "$2a$04$AG864hNeosMGVOZKBePuRejH7ElpHfFBBHTFS6/XFJS4beixwXZB."},
}

func newTestDB(t *testing.T) *DB {
	t.Helper()

	tmpDir := t.TempDir()
	dsn := filepath.Join(tmpDir, "db_test.sqlite?_foreign_keys=on")

	db, err := New(dsn)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := db.Close()
		if err != nil {
			t.Fatal(err)
		}
	})

	err = db.MigrateUp()
	if err != nil {
		t.Fatal(err)
	}

	for _, user := range testUsers {
		id, err := db.InsertUser(user.email, user.hashedPassword)
		if err != nil {
			t.Fatal(err)
		}

		user.id = id
	}

	return db
}
