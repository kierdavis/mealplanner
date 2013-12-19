package mpdb

import (
	"database/sql"
)

const (
	DB_DRIVER = "mysql"
	DB_PARAMS = "?parseDate=true"
	DB_SOURCE = "mp:mp@/mp" + DB_PARAMS
)

type DB struct {
	conn *sql.DB
}

func Connect() (db DB, err error) {
	conn, err := sql.Open(DB_DRIVER, DB_SOURCE)
	if err != nil {
		return DB{}, err
	}
	
	return DB{conn}, nil
}

func (db DB) Close() (err error) {
	return db.conn.Close()
}
