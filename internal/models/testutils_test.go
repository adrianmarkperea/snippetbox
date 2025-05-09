package models

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newTestDB(t *testing.T) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), "postgres://test_web:pass@localhost:5433/test_snippetbox?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		dbpool.Close()
		t.Fatal(err)
	}
	_, err = dbpool.Exec(context.Background(), string(script))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		defer dbpool.Close()

		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = dbpool.Exec(context.Background(), string(script))
		if err != nil {
			t.Fatal(err)
		}
	})

	return dbpool
}
