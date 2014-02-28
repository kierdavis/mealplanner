package mpdb

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"log"
	"time"
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

func InitialiseVersion(q Queryable, debug bool) (err error) {
	var version uint
	err = q.QueryRow("SELECT version FROM version").Scan(&version)
	isNTE := isNonexistentTableError(err)
	
	if err == nil { // All is fine.
		if debug {log.Printf("Version check: OK, current version is %d\n", version)}
		return nil
	
	} else if isNTE || err == sql.ErrNoRows { // No version set.
		if debug {log.Printf("Version check: version not set yet\n")}
		
		if isNTE { // 'version' table does not exist.
			if debug {log.Printf("Version check: creating version table\n")}
			_, err = q.Exec("CREATE TABLE version (version INT UNSIGNED NOT NULL)")
			if err != nil {
				return err
			}
		}
		
		// Check if other tables exist.
		_, err = q.Exec("SELECT meal.id FROM meal LIMIT 1")
		if err == nil { // Table 'meal' exists.
			if debug {log.Printf("Version check: assuming first startup since introduction of versioning\n")}
			version = 0
		
		} else if isNonexistentTableError(err) { // Table 'meal' does not exist.
			if debug {log.Printf("Version check: assuming empty database\n")}
			version = LatestVersion
		
		} else { // Unknown error.
			return err
		}
	
	} else { // Unknown error.
		return err
	}
	
	if debug {log.Printf("Version check: setting version to %d\n", version)}
	_, err = q.Exec("INSERT INTO version VALUES (?)", version)
	return err
}

// ClearTables deletes all records from the entire database.
func ClearTables(q Queryable) (err error) {
	return execList(q, ClearTablesSQLs)
}

// InitDB creates the database tables if they don't exist. If 'debug' is true,
// debug messages are printed. If 'testData' is true, the tables are also
// cleared and test data are added to them.
func InitDB(debug bool, testData bool) (err error) {
	return WithConnection(func(db *sql.DB) (err error) {
		return WithTransaction(db, func(tx *sql.Tx) (err error) {
			err = InitialiseVersion(tx, debug)
			if err != nil {
				return err
			}
			
			err = CreateTables(tx)
			if err != nil {
				return err
			}
			
			err = Migrate(tx, LatestVersion, debug)
			if err != nil {
				return err
			}

			if testData {
				if debug {
					log.Printf("Clearing database and inserting test data.\n")
				}
				
				err = ClearTables(tx)
				if err != nil {
					return err
				}

				err = InsertTestData(tx)
				if err != nil {
					return err
				}
			}

			return nil
		})
	})
}

// InsertTestData inserts some predefined meals and meal plans into the
// database for testing purposes.
func InsertTestData(q Queryable) (err error) {
	err = AddMealWithTags(q, mpdata.MealWithTags{
		Meal: &mpdata.Meal{
			Name:      "Chilli con carne",
			RecipeURL: "http://example.net/chilli",
			Favourite: false,
		},
		Tags: []string{
			"spicy",
			"lentil",
			"rice",
		},
	})

	if err != nil {
		return err
	}

	err = AddMealWithTags(q, mpdata.MealWithTags{
		Meal: &mpdata.Meal{
			Name:      "Carrot and lentil soup",
			RecipeURL: "http://example.net/soup",
			Favourite: false,
		},
		Tags: []string{
			"lentil",
			"soup",
			"quick",
		},
	})

	if err != nil {
		return err
	}

	err = AddMealWithTags(q, mpdata.MealWithTags{
		Meal: &mpdata.Meal{
			Name:      "Nachos",
			RecipeURL: "http://example.net/nachos",
			Favourite: true,
		},
		Tags: []string{
			"spicy",
			"mexican",
		},
	})

	if err != nil {
		return err
	}

	mp1 := &mpdata.MealPlan{
		Notes:     "some notes",
		StartDate: time.Date(2014, time.January, 25, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2014, time.February, 4, 0, 0, 0, 0, time.UTC),
	}

	err = AddMealPlan(q, mp1)
	if err != nil {
		return err
	}

	mp2 := &mpdata.MealPlan{
		Notes:     "some other notes",
		StartDate: time.Date(2014, time.February, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2014, time.February, 8, 0, 0, 0, 0, time.UTC),
	}

	err = AddMealPlan(q, mp2)
	if err != nil {
		return err
	}

	log.Printf("Test meal plans are %d and %d\n", mp1.ID, mp2.ID)

	return nil
}
