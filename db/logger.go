package db

import (
	"database/sql"
	"time"

	"github.com/volatiletech/sqlboiler/boil"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func logQueryTime(logger *log.Logger, startTime time.Time) {
	logger.Debugln("query took " + time.Since(startTime).String())
}

// QueryLogger is used to add extended logging for db queries in debug mode.
type QueryLogger struct {
	DB      *sqlx.DB
	Logger  *log.Logger
	Enabled bool
	Name    string
}

// Query implementation of Executor
func (d *QueryLogger) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if d.Logger != nil {
		if d.Logger != nil {
			d.Logger.Debugln(query)
			d.Logger.Debugln(args...)
			defer logQueryTime(d.Logger, time.Now())
		}
	}
	return d.DB.Query(query, args...)
}

// Exec  implementation of Executor
func (d *QueryLogger) Exec(query string, args ...interface{}) (sql.Result, error) {
	if d.Logger != nil {
		if d.Logger != nil {
			d.Logger.Debugln(query)
			d.Logger.Debugln(args...)
			defer logQueryTime(d.Logger, time.Now())
		}
	}
	return d.DB.Exec(query, args...)
}

// QueryRow  implementation of Executor
func (d *QueryLogger) QueryRow(query string, args ...interface{}) *sql.Row {
	if d.Logger != nil {
		d.Logger.Debugln(query)
		d.Logger.Debugln(args...)
		defer logQueryTime(d.Logger, time.Now())
	}
	return d.DB.QueryRow(query, args...)
}

// Begin  implementation of Executor
func (d *QueryLogger) Begin() (boil.Transactor, error) {
	if d.Logger != nil {
		d.Logger.Debug("->  beginning tx")
	}
	tx, err := d.DB.Begin()
	if err != nil {
		return tx, err
	}
	return &queryLoggerTx{Tx: tx, logger: d.Logger}, nil
}

// Close  implementation of Executor
func (d *QueryLogger) Close() error {
	if d.Logger != nil {
		d.Logger.Printf("closing %s db connection", d.Name)
	}
	return d.DB.Close()
}

type queryLoggerTx struct {
	Tx     *sql.Tx
	logger *log.Logger
}

// Query  implementation of Transaction Executor
func (t *queryLoggerTx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if t.logger != nil {
		t.logger.Debugln("->  " + query)
		defer logQueryTime(t.logger, time.Now())
	}
	return t.Tx.Query(query, args...)
}

// Exec implementation of Transaction Executor
func (t *queryLoggerTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	if t.logger != nil {
		t.logger.Debugln("->  " + query)
		defer logQueryTime(t.logger, time.Now())
	}
	return t.Tx.Exec(query, args...)
}

// QueryRow implementation of Transaction Executor
func (t *queryLoggerTx) QueryRow(query string, args ...interface{}) *sql.Row {
	if t.logger != nil {
		t.logger.Debugln("->  " + query)
		defer logQueryTime(t.logger, time.Now())
	}
	return t.Tx.QueryRow(query, args...)
}

//Commit implementation of Transaction Executor
func (t *queryLoggerTx) Commit() error {
	if t.logger != nil {
		t.logger.Debug("->  committing tx")
	}
	return t.Tx.Commit()
}

// Rollback implementation of Transaction Executor
func (t *queryLoggerTx) Rollback() error {
	if t.logger != nil {
		t.logger.Debug("->  rolling back tx")
	}
	return t.Tx.Rollback()
}
