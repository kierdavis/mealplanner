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
