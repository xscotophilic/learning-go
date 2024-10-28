package models

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	"snippetbox.xscotophilic.art/internal/env"
)

func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", env.DBCreds)
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(strings.TrimSpace(string(script)))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(strings.TrimSpace(string(script)))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	})

	return db
}
