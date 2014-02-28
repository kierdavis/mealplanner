package mpdb

import (
	"fmt"
	"log"
)

type MigrationError struct {
	From    uint
	To      uint
	Message string
}

func (e MigrationError) Error() (msg string) {
	return e.Message
}

type Migration struct {
	From  uint
	To    uint
	Stmts []string
}

func (m *Migration) Apply(q Queryable) (err error) {
	for _, stmt := range m.Stmts {
		_, err = q.Exec(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func FindMigration(from uint, maxTo uint) (m *Migration) {
	var best *Migration

	for _, m = range Migrations {
		if m.From == from && m.To <= maxTo && (best == nil || m.To > best.To) {
			best = m
		}
	}

	return best
}

func GetDatabaseVersion(q Queryable) (v uint, err error) {
	err = q.QueryRow("SELECT version FROM version").Scan(&v)
	return v, err
}

func SetDatabaseVersion(q Queryable, v uint) (err error) {
	_, err = q.Exec("UPDATE version SET version = ?", v)
	return err
}

// Migrate the database from the current version to 'targetVersion'.
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

const LatestVersion = 1

var Migrations = []*Migration{
	// 2014-02-27 - Add 'searchtext' column to 'meal' table.
	&Migration{0, 1, []string{
		"ALTER TABLE meal ADD COLUMN searchtext TEXT NOT NULL",
		"UPDATE meal SET meal.searchtext = " + SearchTextFunc,
	}},
}
