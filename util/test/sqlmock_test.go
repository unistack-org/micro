package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_NewSQLRowsFromFile(t *testing.T) {
	db, c, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows, err := NewSQLRowsFromFile(c, "testdata/Call.csv")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(fmt.Sprintf("%#+v", rows), `cols:[]string{"DepAgrId", "DepAgrNum", "DepAgrDate", "DepAgrCloseDate", "AccCur", "MainFinaccNum", "MainFinaccName", "MainFinaccId", "MainFinaccOpenDt", "DepAgrStatus", "MainFinaccBal", "DepartCode", "CardAccId"}`) {
		t.Fatal("invalid cols after import csv")
	}
}
