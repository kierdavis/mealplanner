package mpdb

import (
	"database/sql"
)

const (
	DB_DRIVER = "mysql"
	DB_PARAMS = "?parseDate=true"
	DB_SOURCE = "mp:mp@/mp" + DB_PARAMS
)

type Queryable interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Prepare(string) (*sql.Stmt, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

func Connect() (db *sql.DB, err error) {
	return sql.Open(DB_DRIVER, DB_SOURCE)
}

type FailedCloseError struct {
	CloseError error
	OriginalError error
}

type WithConnectionFunc func(*sql.DB) error

type WithTransactionFunc func(*sql.Tx) error

func WithConnection(f WithConnectionFunc) (err error) {
	db, err := Connect()
	if err != nil {
		return err
	}
	
	defer func() {
		err2 := db.Close()
		if err2 != nil {
			err = &FailedCloseError{err2, err}
		}
	}()
	
	return f(db)
}

func WithTransaction(db *sql.DB, f WithTransactionFunc) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	
	defer func() {
		var err2 error
		if err != nil {
			err2 = tx.Rollback()
		} else {
			err2 = tx.Commit()
		}
		
		if err2 != nil {
			err = &FailedCloseError{err2, err}
		}
	}()
	
	return f(tx)
}
