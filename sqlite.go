package tt

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	tableTimers       = "timers"
	tableVacationDays = "vacation_days"
)

type DatabaseFilter interface {
	SQL() string
}

type emptyDbFilter struct{}

func (_ emptyDbFilter) SQL() string { return "" }

var EmptyDbFilter emptyDbFilter

type sqlite struct {
	db *sql.DB
}

func (db *sqlite) create() error {
	err := db.createKeyValueTable(tableTimers)
	if err != nil {
		return fmt.Errorf("db: %w", err)
	}
	err = db.createKeyValueTable(tableVacationDays)
	if err != nil {
		return fmt.Errorf("db: %w", err)
	}
	setupStmt := `
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
					  WHERE (json_extract(timers.json, '$.start') <= json_extract(NEW.json, '$.start')
								 AND json_extract(timers.json, '$.stop') > json_extract(NEW.json, '$.start'))
						 OR (json_extract(timers.json, '$.start') > json_extract(NEW.json, '$.start')
								 AND json_extract(timers.json, '$.start') < json_extract(NEW.json, '$.stop')));
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
					  WHERE json_extract(NEW.json, '$.stop') IS NULL
						AND json_extract(timers.json, '$.stop') IS NULL);
	END;`
	_, err = db.db.Exec(setupStmt)
	if err != nil {
		return fmt.Errorf("db: create: %w", err)
	}
	return nil
}

func (db *sqlite) createKeyValueTable(name string) error {
	_, err := db.db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (uuid TEXT PRIMARY KEY,json TEXT NOT NULL);`, name))
	if err != nil {
		return fmt.Errorf("create table: %w", err)
	}
	return nil
}

func (db *sqlite) save(table string, id string, value interface{}) error {
	insertStmt := fmt.Sprintf("INSERT INTO %s VALUES (?, ?);", table)

	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("save: %w: %s", ErrInvalidData, err.Error())
	}

	_, err = db.db.Exec(insertStmt, id, string(b))
	if err != nil {
		return fmt.Errorf("save: %w: %s", ErrInternal, err.Error())
	}
	return nil
}

func (db *sqlite) getOne(table string, filter DatabaseFilter, orderBy OrderBy, target interface{}) error {
	selectStmt := fmt.Sprintf("SELECT `json` FROM %s %s %s;", table, filter.SQL(), orderBy.SQL())

	row := db.db.QueryRow(selectStmt)
	var content string
	err := row.Scan(&content)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get-one: %w", ErrNotFound)
	} else if err != nil {
		return fmt.Errorf("%w: %s", ErrInternal, err.Error())
	}

	err = json.Unmarshal([]byte(content), target)
	if err != nil {
		return fmt.Errorf("get-one: %w: %s", ErrInternal, err.Error())
	}
	return nil
}

func (db *sqlite) getOneById(table string, id string, target interface{}) error {
	selectStmt := fmt.Sprintf("SELECT `json` FROM %s WHERE `uuid` == ?;", table)

	row := db.db.QueryRow(selectStmt, id)
	var content string
	err := row.Scan(&content)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get-one-by-id: %w", ErrNotFound)
	} else if err != nil {
		return fmt.Errorf("%w: %s", ErrInternal, err.Error())
	}

	err = json.Unmarshal([]byte(content), target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInternal, err.Error())
	}
	return nil
}

func (db *sqlite) getMultiple(table string, filter DatabaseFilter, orderBy OrderBy, target interface{}) error {
	fs := filter.SQL()
	os := orderBy.SQL()
	selectStmt := fmt.Sprintf("SELECT `json` FROM %s %s %s;", table, fs, os)

	rows, err := db.db.Query(selectStmt)
	if err != nil {
		return fmt.Errorf("get-multiple: %w: %s", ErrInternal, err.Error())
	}
	var items []string
	for rows.Next() {
		var item string
		if rows.Err() != nil {
			return fmt.Errorf("get-multiple: %w: %s", ErrInternal, rows.Err().Error())
		}
		err = rows.Scan(&item)
		if err != nil {
			return fmt.Errorf("get-multiple: %w: %s", ErrInternal, err.Error())
		}
		items = append(items, item)
	}
	content := fmt.Sprintf("[%s]", strings.Join(items, ","))

	err = json.Unmarshal([]byte(content), target)
	if err != nil {
		return fmt.Errorf("get-multiple: %w: %s", ErrInternal, err.Error())
	}
	return nil
}

func (db *sqlite) update(table string, id string, value interface{}) error {
	updateStmt := fmt.Sprintf("UPDATE %s SET `json` = ? WHERE `uuid` = ?;", table)

	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("update: %w: %s", ErrInvalidData, err.Error())
	}

	res, err := db.db.Exec(updateStmt, string(b), id)
	if err != nil {
		return fmt.Errorf("update: %w: %s", ErrInternal, err.Error())
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update: %w: %s", ErrInternal, err.Error())
	}
	if rowsAffected == 0 {
		return fmt.Errorf("udpate: %w", ErrNotFound)
	}
	return nil
}

func (db *sqlite) remove(table string, id string) error {
	deleteStmt := fmt.Sprintf("DELETE FROM %s WHERE `uuid` = ?;", table)

	res, err := db.db.Exec(deleteStmt, id)
	if err != nil {
		return fmt.Errorf("remove: %w: %s", ErrInternal, err.Error())
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("remove: %w: %s", ErrInternal, err.Error())
	}
	if rowsAffected == 0 {
		return fmt.Errorf("remove: %w", ErrNotFound)
	}
	return nil
}

func (db *sqlite) SaveTimer(timer Timer) error {
	err := timer.Validate()
	if err != nil {
		return err
	}
	return db.save(tableTimers, timer.ID, timer)
}

func (db *sqlite) GetTimer(filter Filter, orderBy OrderBy, timer *Timer) error {
	return db.getOne(tableTimers, filter, orderBy, timer)
}

func (db *sqlite) GetTimerById(id string, timer *Timer) error {
	return db.getOneById(tableTimers, id, timer)
}

func (db *sqlite) GetTimers(filter Filter, orderBy OrderBy, timers *Timers) error {
	return db.getMultiple(tableTimers, filter, orderBy, timers)
}

func (db *sqlite) UpdateTimer(timer Timer) error {
	err := timer.Validate()
	if err != nil {
		return err
	}
	return db.update(tableTimers, timer.ID, timer)
}

func (db *sqlite) RemoveTimer(id string) error {
	return db.remove(tableTimers, id)
}

func (db *sqlite) SaveVacationDay(vacationDay VacationDay) error {
	return db.save(tableVacationDays, vacationDay.ID, vacationDay)
}

func (db *sqlite) GetVacationDay(filter VacationFilter, vacationDay *VacationDay) error {
	return db.getOne(tableVacationDays, filter, OrderBy{}, vacationDay)
}

func (db *sqlite) GetVacationDays(orderBy OrderBy, vacationDays *[]VacationDay) error {
	return db.getMultiple(tableVacationDays, EmptyDbFilter, orderBy, vacationDays)
}

func (db *sqlite) RemoveVacationDay(id string) error {
	return db.remove(tableVacationDays, id)
}

// NewSQLite creates and initializes a new SQLite storage interface. The
// connection is tested using DB.Ping() and the needed tables are created if
// they do not exist.
func NewSQLite(dbFile string) (DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return &sqlite{}, fmt.Errorf("%w: %s", ErrInternal, err.Error())
	}
	err = db.Ping()
	if err != nil {
		return &sqlite{}, fmt.Errorf("%w: %s", ErrInternal, err.Error())
	}
	ttDB := &sqlite{db}
	err = ttDB.create()
	if err != nil {
		return &sqlite{}, fmt.Errorf("unable to init databse: %w", err)
	}
	return ttDB, nil
}
