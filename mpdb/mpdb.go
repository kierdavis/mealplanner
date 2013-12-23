package mpdb

import (
	"database/sql"
	"fmt"
)

const (
	DB_DRIVER = "mysql"
	DB_USER = "mealplanner"
	DB_PASSWORD = "1Ny9IF7WYA6jvSiBXHku"
	DB_ADDRESS = "unix(/var/run/mysqld/mysqld.sock)"
	DB_DATABASE = "mealplanner"
	DB_PARAMS = "parseTime=true"
)

var DB_SOURCE = fmt.Sprintf("%s:%s@%s/%s?%s", DB_USER, DB_PASSWORD, DB_ADDRESS, DB_DATABASE, DB_PARAMS)

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
	What string
	CloseError error
	OriginalError error
}

func (err *FailedCloseError) Error() (msg string) {
	return fmt.Sprintf("%s (when attempting to %s after: %s)", err.CloseError.Error(), err.What, err.OriginalError.Error())
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
			err = &FailedCloseError{"close connection", err2, err}
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
		var what string
		
		if err != nil {
			err2 = tx.Rollback()
			what = "rollback transaction"
		} else {
			err2 = tx.Commit()
			what = "commit transaction"
		}
		
		if err2 != nil {
			err = &FailedCloseError{what, err2, err}
		}
	}()
	
	return f(tx)
}
