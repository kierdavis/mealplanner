// Package mpdb provides routines for manipulating the database whilst
// preserving referential integrity as best as possible.
package mpdb

import (
	"database/sql"
	"fmt"
)

// Database connection constants.
const (
	// The name of the driver to use.
	DB_DRIVER = "mysql"

	// Additional parameters to use.
	DB_PARAMS = "?parseTime=true"
)

var Source = "mealplanner@unix(/var/run/mysqld/mysqld.sock)/mealplanner"

// Interface Queryable represents a type that can be queried (either a *sql.DB
// or *sql.Tx). See documentation on 'database/sql#DB' for information on the
// methods in this interface.
type Queryable interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Prepare(string) (*sql.Stmt, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

type LoggingQueryable struct {
	Q Queryable
}

func (lq LoggingQueryable) Exec(query string, args ...interface{}) (result sql.Result, err error) {
	result, err = lq.Q.Exec(query, args...)
	fmt.Printf("SQL: Exec(%v, %v) -> %v\n", query, args, err)
	return result, err
}

func (lq LoggingQueryable) Prepare(query string) (stmt *sql.Stmt, err error) {
	stmt, err = lq.Q.Prepare(query)
	fmt.Printf("SQL: Prepare(%v) -> %v\n", query, err)
	return stmt, err
}

func (lq LoggingQueryable) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	rows, err = lq.Q.Query(query, args...)
	fmt.Printf("SQL: Query(%v, %v) -> %v\n", query, args, err)
	return rows, err
}

func (lq LoggingQueryable) QueryRow(query string, args ...interface{}) (row *sql.Row) {
	row = lq.Q.QueryRow(query, args...)
	fmt.Printf("SQL: QueryRow(%v, %v) -> %v\n", query, args, row)
	return row
}

// Connect creates a new connection to the database using DB_DRIVER and
// DB_SOURCE.
func Connect() (db *sql.DB, err error) {
	return sql.Open(DB_DRIVER, Source+DB_PARAMS)
}

// Type FailedCloseError contains information regarding a situation where an
// error occurs when closing a resource in response to an earlier error.
type FailedCloseError struct {
	What          string // A string used in the error message to identify what resource was being closed.
	CloseError    error  // The error returned when the resource was closed.
	OriginalError error  // The original error that triggered the closing of the resource.
}

// Error formats the information contained in 'err' into an error message.
func (err *FailedCloseError) Error() (msg string) {
	return fmt.Sprintf("%s\nAdditionally, when attempting to %s: %s", err.OriginalError.Error(), err.What, err.CloseError.Error())
}

// Type WithConnectionFunc represents a function that can be used with
// WithConnection.
type WithConnectionFunc func(*sql.DB) error

// Type WithTransactionFunc represents a function that can be used with
// WithTransaction.
type WithTransactionFunc func(*sql.Tx) error

// WithConnection opens a connection to the database, calls 'f' with the
// database as a parameter, then ensures the database is closed even in the
// event of an error. If an error occurs when closing the database, a
// 'FailedCloseError' is returned.
func WithConnection(f WithConnectionFunc) (err error) {
	// Connect to database
	db, err := Connect()
	if err != nil {
		return err
	}

	// Run the passed function
	err = f(db)

	// Close the database
	closeErr := db.Close()

	// If closing the database caused an error, return a FailedCloseError
	if closeErr != nil {
		err = &FailedCloseError{"close connection", closeErr, err}
	}

	return err
}

// WithTransaction begins a transaction on the given database connection, calls
// 'f' with the transaction as a parameter, then ensures the transaction is
// committed if 'f' completes successfully or rolled back in the event of an
// error. If an error occurs when committing or rolling back the transaction, a
// 'FailedCloseError' is returned.
func WithTransaction(db *sql.DB, f WithTransactionFunc) (err error) {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Run the passed function
	err = f(tx)

	var closeErr error
	var what string

	// Commit or rollback the transaction
	if err != nil {
		closeErr = tx.Rollback()
		what = "roll back transaction"
	} else {
		closeErr = tx.Commit()
		what = "commit transaction"
	}

	// If committing/rolling back the transaction caused an error, return a
	// FailedCloseError
	if closeErr != nil {
		err = &FailedCloseError{what, closeErr, err}
	}

	return err
}
