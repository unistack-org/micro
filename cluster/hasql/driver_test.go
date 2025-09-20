package sql

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.yandex/hasql/v2"
)

func TestDriver(t *testing.T) {
	dbMaster, dbMasterMock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatal(err)
	}
	defer dbMaster.Close()
	dbMasterMock.MatchExpectationsInOrder(false)

	dbMasterMock.ExpectQuery(`.*pg_is_in_recovery.*`).WillReturnRows(
		sqlmock.NewRowsWithColumnDefinition(
			sqlmock.NewColumn("role").OfType("int8", 0),
			sqlmock.NewColumn("replication_lag").OfType("int8", 0)).
			AddRow(1, 0)).
		RowsWillBeClosed().
		WithoutArgs()

	dbMasterMock.ExpectQuery(`SELECT node_name as name`).WillReturnRows(
		sqlmock.NewRows([]string{"name"}).
			AddRow("master-dc1"))

	dbDRMaster, dbDRMasterMock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatal(err)
	}
	defer dbDRMaster.Close()
	dbDRMasterMock.MatchExpectationsInOrder(false)

	dbDRMasterMock.ExpectQuery(`.*pg_is_in_recovery.*`).WillReturnRows(
		sqlmock.NewRowsWithColumnDefinition(
			sqlmock.NewColumn("role").OfType("int8", 0),
			sqlmock.NewColumn("replication_lag").OfType("int8", 0)).
			AddRow(2, 40)).
		RowsWillBeClosed().
		WithoutArgs()

	dbDRMasterMock.ExpectQuery(`SELECT node_name as name`).WillReturnRows(
		sqlmock.NewRows([]string{"name"}).
			AddRow("drmaster1-dc2"))

	dbDRMasterMock.ExpectQuery(`SELECT node_name as name`).WillReturnRows(
		sqlmock.NewRows([]string{"name"}).
			AddRow("drmaster"))

	dbSlaveDC1, dbSlaveDC1Mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatal(err)
	}
	defer dbSlaveDC1.Close()
	dbSlaveDC1Mock.MatchExpectationsInOrder(false)

	dbSlaveDC1Mock.ExpectQuery(`.*pg_is_in_recovery.*`).WillReturnRows(
		sqlmock.NewRowsWithColumnDefinition(
			sqlmock.NewColumn("role").OfType("int8", 0),
			sqlmock.NewColumn("replication_lag").OfType("int8", 0)).
			AddRow(2, 50)).
		RowsWillBeClosed().
		WithoutArgs()

	dbSlaveDC1Mock.ExpectQuery(`SELECT node_name as name`).WillReturnRows(
		sqlmock.NewRows([]string{"name"}).
			AddRow("slave-dc1"))

	dbSlaveDC2, dbSlaveDC2Mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatal(err)
	}
	defer dbSlaveDC2.Close()
	dbSlaveDC1Mock.MatchExpectationsInOrder(false)

	dbSlaveDC2Mock.ExpectQuery(`.*pg_is_in_recovery.*`).WillReturnRows(
		sqlmock.NewRowsWithColumnDefinition(
			sqlmock.NewColumn("role").OfType("int8", 0),
			sqlmock.NewColumn("replication_lag").OfType("int8", 0)).
			AddRow(2, 50)).
		RowsWillBeClosed().
		WithoutArgs()

	dbSlaveDC2Mock.ExpectQuery(`SELECT node_name as name`).WillReturnRows(
		sqlmock.NewRows([]string{"name"}).
			AddRow("slave-dc1"))

	tctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	c, err := NewCluster[Querier](
		WithClusterContext(tctx),
		WithClusterNodeChecker(hasql.PostgreSQLChecker),
		WithClusterNodePicker(NewCustomPicker[Querier](
			CustomPickerMaxLag(100),
		)),
		WithClusterNodes(
			ClusterNode{"slave-dc1", dbSlaveDC1, 1},
			ClusterNode{"master-dc1", dbMaster, 1},
			ClusterNode{"slave-dc2", dbSlaveDC2, 2},
			ClusterNode{"drmaster1-dc2", dbDRMaster, 0},
		),
		WithClusterOptions(
			hasql.WithUpdateInterval[Querier](2*time.Second),
			hasql.WithUpdateTimeout[Querier](1*time.Second),
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	if err = c.WaitForNodes(tctx, hasql.Primary, hasql.Standby); err != nil {
		t.Fatal(err)
	}

	db, err := OpenDBWithCluster(c)
	if err != nil {
		t.Fatal(err)
	}

	// Use context methods
	row := db.QueryRowContext(NodeStateCriterion(t.Context(), hasql.Primary), "SELECT node_name as name")
	if err = row.Err(); err != nil {
		t.Fatal(err)
	}

	nodeName := ""
	if err = row.Scan(&nodeName); err != nil {
		t.Fatal(err)
	}

	if nodeName != "master-dc1" {
		t.Fatalf("invalid node_name %s != %s", "master-dc1", nodeName)
	}
}
