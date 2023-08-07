package test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_SQLFromFile(t *testing.T) {
	ctx := context.TODO()
	db, c, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err = SQLFromFile(c, "testdata/result/01_firstcase/Call_db.csv"); err != nil {
		t.Fatal(err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	rows, err := tx.QueryContext(ctx, "select * from test;")
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var id int64
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			t.Fatal(err)
		}
		if id != 1 || name != "test" {
			t.Fatalf("invalid rows %v %v", id, name)
		}
	}

	if err = rows.Close(); err != nil {
		t.Fatal(err)
	}

	if err = rows.Err(); err != nil {
		t.Fatal(err)
	}

	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}
	if err = c.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func Test_GetCases(t *testing.T) {
	files, err := GetCases("testdata/", nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) == 0 {
		t.Fatalf("no files matching")
	}

	if n := len(files); n != 1 {
		t.Fatalf("invalid number of test cases %d", n)
	}
}
