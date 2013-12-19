package mpdb

import (
	"database/sql"
	"fmt"
	"github.com/kierdavis/mealplanner/mpdata"
	"time"
)

const CreateScoreTableSQL =
	"CREATE TEMPORARY TABLE score ( " +
	"  mealid BIGINT UNSIGNED NOT NULL, " +
	"  score FLOAT NOT NULL " +
	")"

const DropScoreTableSQL =
	"DROP TABLE score"

const FindCsdSQL =
	"SELECT ABS(DATEDIFF(serving.dateserved, ?)) " +
	"FROM serving " +
	"WHERE serving.mealid = ? " +
	"AND serving.dateserved != ? " +
	"ORDER BY ABS(DATEDIFF(serving.dateserved, ?)) ASC " +
	"LIMIT 1"

const InsertScoreSQL =
	"INSERT INTO score " +
	"VALUES (?, ?)"

const ListMealsByScoreSQL =
	"SELECT meal.id, meal.name, meal.recipe, meal.favourite, score.score " +
	"FROM meal " +
	"INNER JOIN score ON score.mealid = meal.id " +
	"ORDER BY score.score DESC"

/*
// A pair consisting of a meal identifier and a corresponding score
type MealScore struct {
	MealID uint64
	Score float32
}
*/

func (db DB) GenerateSuggestions(date time.Time) (suggs []mpdata.MealScore, err error) {
	// Create temporary table
	err = db.createScoreTable()
	if err != nil {
		return nil, err
	}
	
	// Defer dropping the temporary table until this function exits
	defer db.dropScoreTable()
	
	// Get a list of all meals
	meals, err := db.ListMeals(false)
	if err != nil {
		return nil, err
	}
	
	// Prepare findClosestServingDistance query for repeated use
	csdStmt, err := db.conn.Prepare(FindCsdSQL)
	if err != nil {
		return nil, err
	}
	defer csdStmt.Close() // Defer cleanup of the prepared statement
	
	// Prepare insertScore query for repeated use
	insertStmt, err := db.conn.Prepare(InsertScoreSQL)
	if err != nil {
		return nil, err
	}
	defer insertStmt.Close() // Defer cleanup of the prepared statement
	
	for _, meal := range meals {
		// Find closest serving distance
		dist, err := db.findClosestServingDistance(csdStmt, meal.ID, date)
		if err != nil {
			return nil, err
		}
		
		fmt.Println("Suggs tags!!!!!!!")
		
		// Calculate score and insert
		score := mpdata.CalculateScore(meal.Favourite, dist)
		
		err = db.insertScore(insertStmt, meal.ID, score)
		if err != nil {
			return nil, err
		}
		
		//scorePair := MealScore{meal.ID, score}
		//scorePairs = append(scorePairs, scorePair)
		
		/*
		// If batch is full,
		if len(scorePairs) == cap(scorePairs) {
			// Insert batch into score table
			err = db.insertScores(scorePairs)
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
		err = db.insertScores(scorePairs)
		if err != nil {
			return nil, err
		}
	}
	*/
	
	// List all meals, but sorted by score
	suggs, err = db.listMealsByScore()
	if err != nil {
		return nil, err
	}
	
	return suggs, nil
}

func (db DB) createScoreTable() (err error) {
	_, err = db.conn.Exec(CreateScoreTableSQL)
	return err
}

func (db DB) dropScoreTable() (err error) {
	_, err = db.conn.Exec(DropScoreTableSQL)
	return err
}

func (db DB) findClosestServingDistance(stmt *sql.Stmt, mealID uint64, date time.Time) (dist int, err error) {
	err = stmt.QueryRow(date, mealID, date, date).Scan(&dist)
	return dist, err
}

func (db DB) insertScore(stmt *sql.Stmt, mealID uint64, score float32) (err error) {
	_, err = stmt.Exec(mealID, score)
	return err
}

func (db DB) listMealsByScore() (results []mpdata.MealScore, err error) {
	rows, err := db.conn.Query(ListMealsByScoreSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		ms := mpdata.MealScore{
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
