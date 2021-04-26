package storage

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/maxmoehl/tt/test"
	"github.com/maxmoehl/tt/types"

	"github.com/google/uuid"
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

var testSqlite *sqlite

func setupSqliteTest() error {
	err := test.SetConfig(sqliteConfig)
	if err != nil {
		return err
	}
	err = initStorage()
	if err != nil {
		return err
	}
	var ok bool
	testSqlite, ok = s.(*sqlite)
	if !ok {
		return fmt.Errorf("expected storage to be of type *sqlite")
	}
	_, err = testSqlite.db.Exec("DELETE FROM timers WHERE TRUE")
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
	_ = testSqlite.db
}

func TestSqliteGetTimer(t *testing.T) {
	err := setupSqliteTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	timer := types.Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		Stop:    time.Now(),
		Project: "test",
		Task:    "abc",
		Tags:    []string{"a", "b"},
	}
	err = sqliteCreateRecord(timer)
	if err != nil {
		t.Fatal(err.Error())
	}
	sqlTimer, err := s.GetTimer(timer.Uuid)
	if err != nil {
		t.Fatal(err.Error())
	}
	timer.Start = time.Unix(timer.Start.Unix(), 0)
	timer.Stop = time.Unix(timer.Stop.Unix(), 0)
	err = timersEqual(timer, sqlTimer)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestSqliteGetTimerNotExist(t *testing.T) {
	err := setupSqliteTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	// after calling setup the database does not contain any data
	_, err = s.GetTimer(uuid.Must(uuid.NewRandom()))
	if err == nil {
		t.Fatal("expected error because of non-existent timer")
	}
}

func sqliteCreateRecord(timer types.Timer) error {
	insertStmt := `INSERT INTO timers (uuid, start, stop, project, task, tags) VALUES (?, ?, ?, ?, ?, ?)`
	var stop, task, tags interface{}
	if !timer.Stop.IsZero() {
		stop = timer.Stop.Unix()
	}
	if timer.Task != "" {
		task = timer.Task
	}
	if len(timer.Tags) != 0 {
		tags = strings.Join(timer.Tags, ",")
	}
	_, err := testSqlite.db.Exec(insertStmt, timer.Uuid.String(), timer.Start.Unix(), stop, timer.Project, task, tags)
	return err
}

func sqliteRecordExists(timer types.Timer) error {
	selectStmt := `SELECT uuid, start, stop, project, task, tags FROM timers WHERE uuid = ?`
	row := testSqlite.db.QueryRow(selectStmt, timer.Uuid)
	err := row.Err()
	if err != nil {
		return err
	}
	var uuidString, project, task, tagsString *string
	var start, stop *int64
	err = row.Scan(&uuidString, &start, &stop, &project, &task, &tagsString)
	if err != nil {
		return err
	}
	var sqlTimer types.Timer
	if uuidString == nil {
		return fmt.Errorf("found nil uuid")
	}
	sqlTimer.Uuid = uuid.Must(uuid.Parse(*uuidString))
	if start == nil {
		return fmt.Errorf("found nil start time")
	}
	sqlTimer.Start = time.Unix(*start, 0)
	if stop != nil {
		sqlTimer.Stop = time.Unix(*stop, 0)
	}
	if project == nil {
		return fmt.Errorf("found nil project")
	}
	sqlTimer.Project = *project
	if task != nil {
		sqlTimer.Task = *task
	}
	if tagsString != nil {
		sqlTimer.Tags = strings.Split(*tagsString, ",")
	}
	return timersEqual(timer, sqlTimer)
}
