package tt

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func testDb(t *testing.T) DB {
	db, err := NewSQLite(":memory:")
	if err != nil {
		t.Fatalf("unable to create in-memory database: %s", err.Error())
	}
	return db
}

func TestSaveValidTimer(t *testing.T) {
	db := testDb(t)
	err := db.SaveTimer(Timer{
		ID:      uuid.Must(uuid.NewRandom()).String(),
		Start:   time.Now(),
		Stop:    nil,
		Project: "test",
		Task:    "test",
		Tags:    []string{"a", "b"},
	})
	if err != nil {
		t.Fatalf("expected nil error but got '%s'", err.Error())
	}
}

func TestSaveInvalidTimerEmptyProject(t *testing.T) {
	db := testDb(t)
	err := db.SaveTimer(Timer{
		ID:      uuid.Must(uuid.NewRandom()).String(),
		Start:   time.Now(),
		Stop:    nil,
		Project: "",
		Task:    "test",
		Tags:    []string{"a", "b"},
	})
	if err == nil {
		t.Fatal("expected non-nil error but got nil")
	}
	if !errors.Is(err, ErrInvalidTimer) {
		t.Fatalf("expected error to contain '%s', but got '%s'", ErrInvalidTimer, err.Error())
	}
}

func TestGetTimerNotFound(t *testing.T) {
	db := testDb(t)
	var timer Timer
	err := db.GetTimer(EmptyFilter, OrderBy{}, &timer)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected error to contain '%s', but got '%s'", ErrNotFound, err.Error())
	}
}
