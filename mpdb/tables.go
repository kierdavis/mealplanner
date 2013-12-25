package mpdb

import (
	"database/sql"
)

// SQL statements to delete tables.
var DeleteTablesSQLs = []string{
	"DROP TABLE IF EXISTS meal",
	"DROP TABLE IF EXISTS tag",
	"DROP TABLE IF EXISTS mealplan",
	"DROP TABLE IF EXISTS serving",
}

// SQL statements to create tables.
var CreateTablesSQLs = []string{
	"CREATE TABLE IF NOT EXISTS meal ( " +
		"id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, " +
		"name VARCHAR(255) NOT NULL, " +
		"recipe TEXT, " +
		"favourite BOOLEAN NOT NULL, " +
		"PRIMARY KEY (id) " +
	")",
	"CREATE TABLE IF NOT EXISTS tag ( " +
		"mealid BIGINT UNSIGNED NOT NULL, " +
		"tag VARCHAR(64) NOT NULL, " +
		"PRIMARY KEY (mealid, tag) " +
	")",
	"CREATE TABLE IF NOT EXISTS mealplan ( " +
		"id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, " +
		"notes TEXT, " +
		"startdate DATE NOT NULL, " +
		"enddate DATE NOT NULL, " +
		"PRIMARY KEY (id) " +
	")",
	"CREATE TABLE IF NOT EXISTS serving ( " +
		"mealplanid BIGINT UNSIGNED NOT NULL, " +
		"dateserved DATE NOT NULL, " +
		"mealid BIGINT UNSIGNED NOT NULL, " +
		"PRIMARY KEY (mealplanid, dateserved) " +
	")",
}

// SQL statements to clear tables.
var ClearTablesSQLs = []string{
	"DELETE FROM meal",
	"DELETE FROM tag",
	"DELETE FROM mealplan",
	"DELETE FROM serving",
}

// execList runs a list of SQL statements, discarding the results.
func execList(q Queryable, queries []string) (err error) {
	for _, query := range queries {
		_, err = q.Exec(query)
		if err != nil {
			return err
		}
	}
	
	return nil
}

// DeleteTables drops the database tables if they exist.
func DeleteTables(q Queryable) (err error) {
	return execList(q, DeleteTablesSQLs)
} 

// CreateTables creates the database tables if they do not exist.
func CreateTables(q Queryable) (err error) {
	return execList(q, CreateTablesSQLs)
}

// ClearTables deletes all records from the entire database.
func ClearTables(q Queryable) (err error) {
	return execList(q, ClearTablesSQLs)
}

// InitDB creates the database tables if they don't exist. If 'clear' is true,
// the tables are also cleared (in the event that they did exist).
func InitDB(clear bool) (err error) {
	return WithConnection(func(db *sql.DB) (err error) {
		return WithTransaction(db, func(tx *sql.Tx) (err error) {
			err = CreateTables(tx)
			if err != nil {
				return err
			}
			
			if clear {
				err = ClearTables(tx)
				if err != nil {
					return err
				}
			}
			
			return nil
		})
	})
}