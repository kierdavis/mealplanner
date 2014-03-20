package mpdb

import (
	"fmt"
	"log"
)

// A MigrationError is returned when database migration failed.
type MigrationError struct {
	From    uint   // The version we were attempting to migrate from.
	To      uint   // The version we were attempting to migrate to.
	Message string // The error message.
}

func (e MigrationError) Error() (msg string) {
	return e.Message
}

// A Migration represents a possible migration step between two database
// versions.
type Migration struct {
	From  uint     // The version the database must be at before this migration is executed.
	To    uint     // The version the database will be at after this migration is executed.
	Stmts []string // The SQL statements that perform the migration.
}

// Apply runs the migration against a database.
func (m *Migration) Apply(q Queryable) (err error) {
	for _, stmt := range m.Stmts {
		_, err = q.Exec(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

// FindMigration finds a migration step that starts at version 'from' and
// finishes anywhere up to and including version 'maxTo'. If there are multiple
// possible choices, the one that spans the most versions (i.e. the one with
// the highest "to" version) is returned.
func FindMigration(from uint, maxTo uint) (m *Migration) {
	var best *Migration

	for _, m = range Migrations {
		if m.From == from && m.To <= maxTo && (best == nil || m.To > best.To) {
			best = m
		}
	}

	return best
}

// GetDatabaseVersion fetches and returns the version number of the database.
func GetDatabaseVersion(q Queryable) (v uint, err error) {
	err = q.QueryRow("SELECT version FROM version").Scan(&v)
	return v, err
}

// SetDatabaseVersion updates the version number in the database.
func SetDatabaseVersion(q Queryable, v uint) (err error) {
	_, err = q.Exec("UPDATE version SET version = ?", v)
	return err
}

// Migrate migrates the database from the current version to 'targetVersion'.
// If 'debug' is true, messages are printed to stdout describing the operations
// taking place.
func Migrate(q Queryable, targetVersion uint, debug bool) (err error) {
	currentVersion, err := GetDatabaseVersion(q)
	if err != nil {
		return err
	}

	if currentVersion > targetVersion {
		return MigrationError{
			From:    currentVersion,
			To:      targetVersion,
			Message: fmt.Sprintf("Cannot migrate to an earlier version of the database (%d) from the current version (%d)", targetVersion, currentVersion),
		}
	}

	if debug {
		log.Printf("Migration: Database is at version %d, migration target is %d. Checking for available migrations.\n", currentVersion, targetVersion)
	}

	for currentVersion < targetVersion {
		m := FindMigration(currentVersion, targetVersion)
		if m == nil {
			return MigrationError{
				From:    currentVersion,
				To:      targetVersion,
				Message: fmt.Sprintf("No migration defined between versions %d and %d", currentVersion, targetVersion),
			}
		}

		if debug {
			log.Printf("Migration: Executing migration from version %d to %d.\n", m.From, m.To)
		}

		err = m.Apply(q)
		if err != nil {
			return err
		}

		currentVersion = m.To
		err = SetDatabaseVersion(q, currentVersion)
		if err != nil {
			return err
		}
	}

	if debug {
		log.Printf("Migration: Done. Database is now at version %d.\n", currentVersion)
	}

	return nil
}

// LatestVersion is the latest version a database could be in.
const LatestVersion = 1

// Migrations contains a list of all possible migration steps.
var Migrations = []*Migration{
	// 2014-02-27 - Add 'searchtext' column to 'meal' table.
	&Migration{0, 1, []string{
		"ALTER TABLE meal ADD COLUMN searchtext TEXT NOT NULL",
		"UPDATE meal SET meal.searchtext = " + SearchTextExpr,
	}},
}
