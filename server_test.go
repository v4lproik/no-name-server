package main

import (
	"testing"
	"github.com/yinkozi/no-name-domain"
	"database/sql"
	"os"
)

const dbpath = "foo_test.db"

var db *sql.DB

func TestMain(t *testing.M) {
	db = InitDB(dbpath)
	defer db.Close()

	CreateTable(db)

	os.Exit(t.Run())
}

func TestInsertingTwoItemsInDatabaseShouldReturnTwoItems(t *testing.T) {
	// given
	items := []domain.Report{
		domain.Report{"1", "TEST", nil},
		domain.Report{"2", "TEST", nil},
	}
	StoreItem(db, items[0])
	StoreItem(db, items[1])

	// when
	readAllItems := ReadAllItems(db)

	// then
	if len(readAllItems) != 2 {
		t.Errorf("Two items should have been stored in the database")
	}
}

func TestDeletingTwoItemsInDatabaseShouldReturnNoMoreItems(t *testing.T) {
	items := []domain.Report{
		domain.Report{"1", "TEST", nil},
		domain.Report{"2", "TEST", nil},
	}
	StoreItem(db, items[0])
	StoreItem(db, items[1])

	//when
	DeleteItem(db, items[0].Id)
	DeleteItem(db, items[1].Id)
	readAllItems := ReadAllItems(db)

	// then
	if len(readAllItems) > 0 {
		t.Errorf("The two items should have been removed from the database")
	}
}

func TestUpdateTwoItemsInDatabaseShouldReturnTwoUpdatedItems(t *testing.T) {
	items := []domain.Report{
		domain.Report{"1", "TEST", nil},
		domain.Report{"2", "TEST", nil},
	}
	StoreItem(db, items[0])
	StoreItem(db, items[1])

	//when
	items[0].TypeEnum = "TEST1"
	UpdateItem(db, items[0])
	items[1].TypeEnum = "TEST2"
	UpdateItem(db, items[1])

	readAllItems := ReadAllItems(db)

	// then
	if readAllItems[0].TypeEnum != "TEST1" || readAllItems[1].TypeEnum != "TEST2"  {
		t.Errorf("The two items should have been updated from the database")
	}
}