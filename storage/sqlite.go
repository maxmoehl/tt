package storage

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/maxmoehl/tt/config"
	"github.com/maxmoehl/tt/types"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const (
	sqlOperatorLike   = "LIKE"
	sqlOperatorEquals = "="
)

type scanable interface {
	Scan(dest ...interface{}) error
}

type sqlite struct {
	db *sql.DB
}

func (db *sqlite) GetTimer(uuid uuid.UUID) (types.Timer, error) {
	selectStmt := `
		SELECT *
		FROM timers
		WHERE timers.uuid = ?;`
	row := db.db.QueryRow(selectStmt, uuid.String())
	return db.scanRow(row)
}

func (db *sqlite) GetRunningTimer() (types.Timer, error) {
	selectStmt := `
		SELECT *
		FROM timers
		WHERE stop ISNULL;`
	row := db.db.QueryRow(selectStmt)
	return db.scanRow(row)
}

func (db *sqlite) GetLastTimer(running bool) (types.Timer, error) {
	selectStmt := `
		SELECT *
		FROM timers
		ORDER BY start DESC
		LIMIT 2;`
	rows, err := db.db.Query(selectStmt)
	if err != nil {
		return types.Timer{}, err
	}
	var timers types.Timers
	for rows.Next() {
		t, err := db.scanRow(rows)
		if err != nil {
			return types.Timer{}, err
		}
		timers = append(timers, t)
	}
	if rows.Err() != nil {
		return types.Timer{}, rows.Err()
	}
	if timers == nil {
		return types.Timer{}, fmt.Errorf("no timers found")
	}
	timer := timers.Last(running)
	if timer.IsZero() {
		return types.Timer{}, fmt.Errorf("no timer found")
	}
	return timer, nil
}

func (db *sqlite) GetTimers(filter types.Filter) (types.Timers, error) {
	selectStmt := fmt.Sprintf(`
		SELECT *
		FROM timers
		WHERE %s;`, getWhereClause(filter))
	rows, err := db.db.Query(selectStmt)
	if err != nil {
		return types.Timers{}, err
	}
	var timers types.Timers
	var t types.Timer
	for rows.Next() {
		t, err = db.scanRow(rows)
		if err != nil {
			return nil, err
		}
		// We still filter the timers since the sql filtering for tags
		// is questionable at best
		if filter.Match(t) {
			timers = append(timers, t)
		}
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return timers, nil
}

func (db *sqlite) RunningTimerExists() (bool, error) {
	selectStmt := `
		SELECT *
		FROM timers
		WHERE stop ISNULL;`
	rows, err := db.db.Query(selectStmt)
	if err != nil {
		return false, err
	}
	found := rows.Next()
	if !found {
		err = rows.Err()
		if err != nil {
			return false, err
		}
	}
	return found, nil
}

func (db *sqlite) StoreTimer(timer types.Timer) error {
	insertStmt := `
		INSERT INTO timers (uuid, start, stop, project, task, tags)
		VALUES (?, ?, ?, ?, ?, ?);`
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
	return err
}

func (db *sqlite) UpdateTimer(timer types.Timer) error {
	updateStmt := `
		UPDATE timers
		SET stop = ?
		WHERE uuid = ?;`
	res, err := db.db.Exec(updateStmt, timer.Stop.Unix(), timer.Uuid.String())
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("timer does not exist")
	}
	return nil
}

func (db *sqlite) create() error {
	createStmt := `
		CREATE TABLE IF NOT EXISTS
			timers
		(
			uuid    TEXT PRIMARY KEY,
			start   INTEGER NOT NULL,
			stop    INTEGER,
			project TEXT NOT NULL,
			task    TEXT,
			tags    TEXT
		);`
	_, err := db.db.Exec(createStmt)
	return err
}

func (db *sqlite) scanRow(row scanable) (types.Timer, error) {
	var id, project, task, tags *string
	var start, stop *int64
	err := row.Scan(&id, &start, &stop, &project, &task, &tags)
	if err != nil {
		return types.Timer{}, err
	}
	if id == nil {
		return types.Timer{}, fmt.Errorf("found nil uuid")
	}
	UUID, err := uuid.Parse(*id)
	if err != nil {
		return types.Timer{}, err
	}
	if start == nil {
		return types.Timer{}, fmt.Errorf("found nil start time")
	}
	startTime := time.Unix(*start, 0)
	var stopTime time.Time
	if stop != nil {
		stopTime = time.Unix(*stop, 0)
	}
	if project == nil {
		return types.Timer{}, fmt.Errorf("found nil project string")
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
	return types.Timer{
		Uuid:    UUID,
		Start:   startTime,
		Stop:    stopTime,
		Project: *project,
		Task:    taskString,
		Tags:    tagList,
	}, nil
}

func NewSQLite() (types.Storage, error) {
	storage := &sqlite{}
	var err error
	storage.db, err = sql.Open("sqlite3", filepath.Join(config.HomeDir(), "storage.db"))
	if err != nil {
		return nil, err
	}
	err = storage.db.Ping()
	if err != nil {
		return nil, err
	}
	err = storage.create()
	if err != nil {
		return nil, err
	}
	return storage, nil
}

func getWhereClause(f types.Filter) string {
	var filters []string
	projects := convertFilter("project", f.Project, sqlOperatorEquals)
	if projects != "" {
		filters = append(filters, projects)
	}
	tasks := convertFilter("task", f.Task, sqlOperatorEquals)
	if tasks != "" {
		filters = append(filters, tasks)
	}
	tags := convertFilter("tags", f.Tags, sqlOperatorLike)
	if tags != "" {
		filters = append(filters, tags)
	}
	if !f.Since.IsZero() {
		filters = append(filters, fmt.Sprintf("start > %d", f.Since.Unix()))
	}
	if !f.Until.IsZero() {
		filters = append(filters, fmt.Sprintf("stop < %d", f.Until.Unix()))
	}
	// if there are no filters return TRUE to match all values
	if len(filters) == 0 {
		return "TRUE"
	}
	return strings.Join(filters, " AND ")
}

func convertFilter(key string, values []string, operator string) string {
	if len(values) == 0 {
		return ""
	}
	b := strings.Builder{}
	for i, v := range values {
		if i > 0 {
			b.WriteString(" OR ")
		} else {
			b.WriteString("(")
		}
		switch operator {
		case sqlOperatorEquals:
			b.WriteString(key)
			b.WriteString("='")
			b.WriteString(v)
			b.WriteString("'")
		case sqlOperatorLike:
			b.WriteString(key)
			b.WriteString(" LIKE '%")
			b.WriteString(v)
			b.WriteString("%'")
		}
	}
	b.WriteString(")")
	return b.String()
}