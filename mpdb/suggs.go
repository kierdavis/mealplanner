package mpdb

import (
	"database/sql"
	"fmt"
	"github.com/kierdavis/mealplanner/mpdata"
	"time"
)

// SQL to create the temporary score table.
const CreateScoreTableSQL = "CREATE TEMPORARY TABLE score ( " +
	"  mealid BIGINT UNSIGNED NOT NULL, " +
	"  score FLOAT NOT NULL, " +
	"  PRIMARY KEY (mealid) " +
	")"

// SQL to drop the temporary score table.
const DropScoreTableSQL = "DROP TABLE score"

// SQL to find the closest serving distance to a given date.
const FindCsdSQL = "SELECT ABS(DATEDIFF(serving.dateserved, ?)) " +
	"FROM serving " +
	"WHERE serving.mealid = ? " +
	"AND serving.dateserved != ? " +
	"ORDER BY ABS(DATEDIFF(serving.dateserved, ?)) ASC " +
	"LIMIT 1"

// SQL to insert a meal identifier and score into the temporary score table.
const InsertScoreSQL = "INSERT INTO score VALUES (?, ?)"

// SQL to list meals sorted by score.
const ListMealsByScoreSQL = "SELECT meal.id, meal.name, meal.recipe, meal.favourite, score.score " +
	"FROM meal " +
	"INNER JOIN score ON score.mealid = meal.id " +
	"ORDER BY score.score DESC"

// GenerateSuggestions takes a date to generate suggestions for, and produces
// a list of pairs of meals and scores.
func GenerateSuggestions(q Queryable, date time.Time) (suggs []mpdata.MealWithScore, err error) {
	// Create temporary table
	err = createScoreTable(q)
	if err != nil {
		return nil, err
	}

	// Defer dropping the temporary table until this function exits
	defer dropScoreTable(q)

	// Get a list of all meals
	meals, err := ListMeals(q, false)
	if err != nil {
		return nil, err
	}

	// Prepare findClosestServingDistance query for repeated use
	csdStmt, err := q.Prepare(FindCsdSQL)
	if err != nil {
		return nil, err
	}
	defer csdStmt.Close() // Defer cleanup of the prepared statement

	// Prepare insertScore query for repeated use
	insertStmt, err := q.Prepare(InsertScoreSQL)
	if err != nil {
		return nil, err
	}
	defer insertStmt.Close() // Defer cleanup of the prepared statement

	for _, meal := range meals {
		// Find closest serving distance
		dist, err := findClosestServingDistance(q, csdStmt, meal.ID, date)
		if err != nil {
			return nil, err
		}

		fmt.Println("Suggs tags!!!!!!!")

		// Calculate score and insert
		score := mpdata.CalculateScore(meal.Favourite, dist)

		err = insertScore(q, insertStmt, meal.ID, score)
		if err != nil {
			return nil, err
		}

		//scorePair := MealWithScore{meal.ID, score}
		//scorePairs = append(scorePairs, scorePair)

		/*
			// If batch is full,
			if len(scorePairs) == cap(scorePairs) {
				// Insert batch into score table
				err = insertScores(scorePairs)
				if err != nil {
					return nil, err
				}

				// Truncate batch buffer to be empty
				scorePairs = scorePairs[:0]
			}
		*/
	}

	/*
		// If there are scores not yet inserted,
		if len(scorePairs) > 0 {
			// Insert them
			err = insertScores(scorePairs)
			if err != nil {
				return nil, err
			}
		}
	*/

	// List all meals, but sorted by score
	suggs, err = listMealsByScore(q)
	if err != nil {
		return nil, err
	}

	return suggs, nil
}

// createScoreTable creates the temporary score table.
func createScoreTable(q Queryable) (err error) {
	_, err = q.Exec(CreateScoreTableSQL)
	return err
}

// dropScoreTable drops the temporary score table.
func dropScoreTable(q Queryable) (err error) {
	_, err = q.Exec(DropScoreTableSQL)
	return err
}

// findClosestServingDistance finds the distance from 'date' to the closest
// serving of the meal identified by 'mealID'. 'stmt' should be a prepared
// statement compiled from FindCsdSQL.
func findClosestServingDistance(q Queryable, stmt *sql.Stmt, mealID uint64, date time.Time) (dist int, err error) {
	err = stmt.QueryRow(date, mealID, date, date).Scan(&dist)
	if err == sql.ErrNoRows {
		return -1, nil
	}
	return dist, err
}

// insertScore inserts 'mealID' and 'score' into a new record in the score
// table. 'stmt' should be a prepared statement compiled from InsertScoreSQL.
func insertScore(q Queryable, stmt *sql.Stmt, mealID uint64, score float32) (err error) {
	_, err = stmt.Exec(mealID, score)
	return err
}

// listMealsByScore returns a list of pairs of meals and their scores, sorted
// by score.
func listMealsByScore(q Queryable) (results []mpdata.MealWithScore, err error) {
	rows, err := q.Query(ListMealsByScoreSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		ms := mpdata.MealWithScore{
			Meal: &mpdata.Meal{},
		}

		err = rows.Scan(&ms.Meal.ID, &ms.Meal.Name, &ms.Meal.RecipeURL, &ms.Meal.Favourite, &ms.Score)
		if err != nil {
			return nil, err
		}

		results = append(results, ms)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}
