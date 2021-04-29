package storage

import (
	"testing"

	"github.com/maxmoehl/tt/test"
)

const sqliteConfig = `precision: m
workHours: 8
storageType: sqlite
workDays:
  monday: true
  tuesday: true
  wednesday: true
  thursday: true
  friday: true
  saturday: false
  sunday: false`

func setupSqliteTest() error {
	err := test.SetConfig(sqliteConfig)
	if err != nil {
		return err
	}
	err = initStorage()
	if err != nil {
		return err
	}
	return nil
}

func TestNewSQLite(t *testing.T) {
	err := setupSqliteTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	s, err := NewSQLite()
	if err != nil {
		t.Fatal(err.Error())
	}
	_, ok := s.(*sqlite)
	if !ok {
		t.Fatal("expected storage to be of type *sqlite")
	}
}
