// Package mpdb provides routines for manipulating the database whilst
// preserving referential integrity as best as possible.
package mpdb

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
)

// DBDriver is the driver name used when connecting to the database.
const DBDriver = "mysql"

// DBParams are extra parameters required for the database routines to function.
const DBParams = "?parseTime=true"

// DBSource identifies how to connect to the database. It should take the form
// "USER:PASS@unix(/PATH/TO/SOCKET)/DBNAME" or "USER:PASS@tcp(HOST:PORT)/DBNAME".
// By default, it will attempt to connect via the local Unix socket to the
// 'mealplanner' database, with username 'mealplanner' and no password.
var DBSource = "mealplanner@unix(/var/run/mysqld/mysqld.sock)/mealplanner"

// Queryable represents a type that can be queried (either a *sql.DB
// or *sql.Tx). See documentation on 'database/sql#DB' for information on the
// methods in this interface.
type Queryable interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Prepare(string) (*sql.Stmt, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

// LoggingQueryable wraps a Queryable while logging all executions of its
// functions to standard output. It is intended for debugging purposes.
type LoggingQueryable struct {
	Q Queryable
}

// Exec executes a query without returning any rows. The args are for any
// placeholder parameters in the query.
func (lq LoggingQueryable) Exec(query string, args ...interface{}) (result sql.Result, err error) {
	result, err = lq.Q.Exec(query, args...)
	log.Printf("SQL: Exec(%v, %v) -> %v\n", query, args, err)
	return result, err
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned
// statement.
func (lq LoggingQueryable) Prepare(query string) (stmt *sql.Stmt, err error) {
	stmt, err = lq.Q.Prepare(query)
	log.Printf("SQL: Prepare(%v) -> %v\n", query, err)
	return stmt, err
}

// Query executes a query that returns rows, typically a SELECT. The args are
// for any placeholder parameters in the query.
func (lq LoggingQueryable) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	rows, err = lq.Q.Query(query, args...)
	log.Printf("SQL: Query(%v, %v) -> %v\n", query, args, err)
	return rows, err
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always return a non-nil value. Errors are deferred until Row's Scan
// method is called.
func (lq LoggingQueryable) QueryRow(query string, args ...interface{}) (row *sql.Row) {
	row = lq.Q.QueryRow(query, args...)
	log.Printf("SQL: QueryRow(%v, %v) -> %v\n", query, args, row)
	return row
}

// Connect creates a new connection to the database using DBDriver and
// DB_SOURCE.
func Connect() (db *sql.DB, err error) {
	return sql.Open(DBDriver, DBSource+DBParams)
}

// FailedCloseError contains information regarding a situation where an error
// occurs when closing a resource in response to an earlier error.
type FailedCloseError struct {
	What          string // A string used in the error message to identify what resource was being closed.
	CloseError    error  // The error returned when the resource was closed.
	OriginalError error  // The original error that triggered the closing of the resource.
}

// Error formats the information contained in 'err' into an error message.
func (err *FailedCloseError) Error() (msg string) {
	return fmt.Sprintf("%s\nAdditionally, when attempting to %s: %s", err.OriginalError.Error(), err.What, err.CloseError.Error())
}

// WithConnectionFunc represents a function that can be used with
// WithConnection.
type WithConnectionFunc func(*sql.DB) error

// WithTransactionFunc represents a function that can be used with
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

func isNonexistentTableError(err error) bool {
	mysqlError, isMysqlError := err.(*mysql.MySQLError)
	return isMysqlError && mysqlError.Number == 1146
}
