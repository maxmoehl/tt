package tt

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // import database driver
)

type scanable interface {
	Scan(dest ...interface{}) error
}

type sqlite struct {
	db *sql.DB
}

func (db *sqlite) GetTimer(uuid uuid.UUID) (Timer, Error) {
	selectStmt := `
		SELECT *
		FROM timers
		WHERE timers.uuid = ?;`
	row := db.db.QueryRow(selectStmt, uuid.String())
	timer, err := db.scanRow(row)
	if err != nil {
		return Timer{}, NewError("unable to get timer").WithCause(err)
	}
	return timer, nil
}

func (db *sqlite) GetLastTimer(running bool) (Timer, Error) {
	selectStmt := `
		SELECT *
		FROM timers
		ORDER BY start DESC
		LIMIT 2;`
	rows, err := db.db.Query(selectStmt)
	if err != nil {
		return Timer{}, ErrInternalError.WithCause(err)
	}
	var timers Timers
	for rows.Next() {
		t, err := db.scanRow(rows)
		if err != nil {
			return Timer{}, NewError("unable to get last timer").WithCause(err)
		}
		timers = append(timers, t)
	}
	if rows.Err() != nil {
		return Timer{}, ErrInternalError.WithCause(rows.Err())
	}
	if len(timers) == 0 {
		return Timer{}, ErrNotFound
	}
	timer := timers.Last(running)
	if timer.IsZero() {
		return Timer{}, ErrNotFound
	}
	return timer, nil
}

func (db *sqlite) GetTimers(filter Filter) (Timers, Error) {
	selectStmt := fmt.Sprintf(`
		SELECT *
		FROM timers
		WHERE %s;`, filter.SQL())
	rows, err := db.db.Query(selectStmt)
	if err != nil {
		return nil, ErrInternalError.WithCause(err)
	}
	var timers Timers
	var t Timer
	for rows.Next() {
		t, err = db.scanRow(rows)
		if err != nil {
			return nil, NewError("unable to get timers").WithCause(err)
		}
		// We still filter the timers since the sql filtering for tags
		// is questionable at best
		if filter.Match(t) {
			timers = append(timers, t)
		}
	}
	if rows.Err() != nil {
		return nil, ErrInternalError.WithCause(rows.Err())
	}
	return timers, nil
}

func (db *sqlite) StoreTimer(timer Timer) Error {
	insertStmt := `
		INSERT INTO timers (uuid, start, stop, project, task, tags)
		VALUES (?, ?, ?, ?, ?, ?);`
	if timer.IsZero() {
		return ErrInvalidData.WithCause(NewErrorf("timer is zero"))
	}
	id := timer.Uuid.String()
	start := timer.Start.Unix()
	var stop, task, tags interface{}
	if !timer.Stop.IsZero() {
		stop = timer.Stop.Unix()
	}
	if timer.Task != "" {
		task = timer.Task
	}
	if len(timer.Tags) > 0 {
		tags = strings.Join(timer.Tags, ",")
	}
	_, err := db.db.Exec(insertStmt, id, start, stop, timer.Project, task, tags)
	if err != nil {
		return ErrInternalError.WithCause(err)
	}
	return nil
}

func (db *sqlite) UpdateTimer(timer Timer) Error {
	updateStmt := `
		UPDATE timers
		SET stop = ?
		WHERE uuid = ? AND stop IS NULL;`
	res, err := db.db.Exec(updateStmt, timer.Stop.Unix(), timer.Uuid.String())
	if err != nil {
		return ErrInternalError.WithCause(err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return ErrInternalError.WithCause(err)
	}
	if rowsAffected == 0 {
		return ErrNotFound.WithCause(NewError("either the timer could not be found, or it already has a stop time"))
	}
	return nil
}

func (db *sqlite) create() Error {
	setupStmt := `
		-- create the tables
		CREATE TABLE IF NOT EXISTS
			timers
		(
			uuid    TEXT PRIMARY KEY,
			start   INTEGER NOT NULL,
			stop    INTEGER,
			project TEXT    NOT NULL,
			task    TEXT,
			tags    TEXT
		);
		-- create trigger to prevent collisions
		CREATE TRIGGER IF NOT EXISTS noCollisions
			BEFORE INSERT
			ON timers
			FOR EACH ROW
		BEGIN
			SELECT RAISE(ROLLBACK, 'new timer collides with existing one')
			WHERE EXISTS(
						  SELECT 1
						  FROM timers
						  WHERE (timers.start <= NEW.start AND timers.stop > NEW.start)
							 OR (timers.start > NEW.start AND timers.start < NEW.stop));
		END;
		-- create trigger to prevent multiple running timers
		CREATE TRIGGER IF NOT EXISTS onlyOneRunning
			BEFORE INSERT
			ON timers
			FOR EACH ROW
		BEGIN
			SELECT RAISE(ROLLBACK, 'running timer already exists, cannot have two running timers')
			WHERE EXISTS(
						  SELECT 1
						  FROM timers
						  WHERE NEW.stop IS NULL
							AND timers.stop IS NULL);
		END;`
	_, err := db.db.Exec(setupStmt)
	if err != nil {
		return ErrInternalError.WithCause(err)
	}
	return nil
}

func (db *sqlite) scanRow(row scanable) (Timer, Error) {
	var id, project, task, tags *string
	var start, stop *int64
	err := row.Scan(&id, &start, &stop, &project, &task, &tags)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Timer{}, ErrNotFound
		}
		return Timer{}, ErrInternalError.WithCause(err)
	}
	if id == nil {
		return Timer{}, ErrInternalError.WithCause(NewError("found nil uuid"))
	}
	UUID, err := uuid.Parse(*id)
	if err != nil {
		return Timer{}, ErrInternalError.WithCause(err)
	}
	if start == nil {
		return Timer{}, ErrInternalError.WithCause(NewError("found nil start time"))
	}
	var stopTime time.Time
	if stop != nil {
		stopTime = time.Unix(*stop, 0)
	}
	if project == nil {
		return Timer{}, ErrInternalError.WithCause(NewError("found nil project string"))
	}
	var taskString string
	if task != nil {
		taskString = *task
	} else {
		taskString = ""
	}
	var tagList []string
	if tags != nil {
		tagList = strings.Split(*tags, ",")
	}
	return Timer{
		Uuid:    UUID,
		Start:   time.Unix(*start, 0),
		Stop:    stopTime,
		Project: *project,
		Task:    taskString,
		Tags:    tagList,
	}, nil
}

// NewSQLite creates and initializes a new SQLite storage interface.
// The connection is tested using sql.DB.Ping() and the timers' table
// is created if it does not exist.
func NewSQLite(c Config) (Storage, Error) {
	storage := &sqlite{}
	var e error
	storage.db, e = sql.Open("sqlite3", filepath.Join(c.HomeDir(), "storage.db"))
	if e != nil {
		return nil, ErrInternalError.WithCause(e)
	}
	e = storage.db.Ping()
	if e != nil {
		return nil, ErrInternalError.WithCause(e)
	}
	err := storage.create()
	if err != nil {
		return nil, NewError("unable to create table (if it doesn't exist yet) db").WithCause(err)
	}
	return storage, nil
}
