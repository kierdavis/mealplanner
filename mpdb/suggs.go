package mpdb

import (
	"github.com/kierdavis/mealplanner/mpdata"
	"time"
)

const GetTagDistsSQL = "SELECT tag.tag, ABS(DATEDIFF(serving.dateserved, ?)) " +
	"FROM tag " +
	"INNER JOIN serving ON serving.mealid = tag.mealid"

const ListMealsForSuggsSQL = "SELECT meal.id, meal.favourite, MIN(ABS(DATEDIFF(serving.dateserved, ?))) " +
	"FROM meal " +
	"LEFT JOIN serving ON meal.id = serving.mealid " +
	"GROUP BY meal.id"

// SQL to create the temporary score table.
const CreateScoreTableSQL = "CREATE TEMPORARY TABLE score ( " +
	"  mealid BIGINT UNSIGNED NOT NULL, " +
	"  score FLOAT NOT NULL, " +
	"  PRIMARY KEY (mealid) " +
	")"

// SQL to insert a meal identifier and score into the temporary score table.
const InsertScoreSQL = "INSERT INTO score VALUES (?, ?)"

// SQL to drop the temporary score table.
const DropScoreTableSQL = "DROP TABLE score"

// SQL to list meals sorted by score.
const ListMealsByScoreSQL = "SELECT meal.id, meal.name, meal.recipe, meal.favourite, score.score " +
	"FROM meal " +
	"INNER JOIN score ON meal.id = score.mealid " +
	"ORDER BY score.score DESC"

// GenerateSuggestions takes a date to generate suggestions for, and produces
// a list of pairs of meals and scores.
func GenerateSuggestions(q Queryable, date time.Time) (suggs []mpdata.MealWithScore, err error) {
	s := mpdata.NewScorer()
	
	err = calculateTagScores(q, s, date)
	if err != nil {
		return nil, err
	}
	
	_, err = q.Exec(CreateScoreTableSQL)
	if err != nil {
		return nil, err
	}
	
	err = scoreMeals(q, s, date)
	if err != nil {
		return nil, err
	}
	
	suggs, err = readScoreTable(q)
	if err != nil {
		return nil, err
	}
	
	_, err = q.Exec(DropScoreTableSQL)
	if err != nil {
		return nil, err
	}
	
	return suggs, nil
}

func calculateTagScores(q Queryable, s *mpdata.Scorer, date time.Time) (err error) {
	rows, err := q.Query(GetTagDistsSQL, date)
	if err != nil {
		return err
	}
	defer rows.Close()
	
	var tag string
	var dist int
	
	for rows.Next() {
		err = rows.Scan(&tag, &dist)
		if err != nil {
			return err
		}
		
		s.AddTagScore(tag, dist)
	}
	
	err = rows.Err()
	if err != nil {
		return err
	}
	
	return nil
}

func scoreMeals(q Queryable, s *mpdata.Scorer, date time.Time) (err error) {
	getTagsStmt, err := q.Prepare(GetMealTagsSQL)
	if err != nil {
		return err
	}
	defer getTagsStmt.Close()
	
	insertScoreStmt, err := q.Prepare(InsertScoreSQL)
	if err != nil {
		return err
	}
	defer insertScoreStmt.Close()

	rows, err := q.Query(ListMealsForSuggsSQL)
	if err != nil {
		return err
	}
	defer rows.Close()
	
	var mealID uint64
	var favourite bool
	var csd int
	var tags []string
	
	for rows.Next() {
		err = rows.Scan(&mealID, &favourite, &csd)
		if err != nil {
			return err
		}
		
		tags, err = getMealTagsPrepared(getTagsStmt, mealID)
		if err != nil {
			return err
		}
		
		score := s.ScoreMeal(favourite, csd, tags)
		_, err = insertScoreStmt.Exec(mealID, score)
		if err != nil {
			return err
		}
	}
	
	return nil
}

// readScoreTable returns a list of pairs of meals and their scores, sorted
// by score.
func readScoreTable(q Queryable) (results []mpdata.MealWithScore, err error) {
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














/*




// SQL to find the closest serving distance to a given date.
const FindCsdSQL = "SELECT ABS(DATEDIFF(serving.dateserved, ?)) " +
	"FROM serving " +
	"WHERE serving.mealid = ? " +
	"AND serving.dateserved != ? " +
	"ORDER BY ABS(DATEDIFF(serving.dateserved, ?)) ASC " +
	"LIMIT 1"

// SQL to insert a meal identifier and score into the temporary score table.
const InsertScoreSQL = "INSERT INTO score VALUES (?, ?)"
	
	

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
	}

		// If there are scores not yet inserted,
		if len(scorePairs) > 0 {
			// Insert them
			err = insertScores(scorePairs)
			if err != nil {
				return nil, err
			}
		}

	// List all meals, but sorted by score
	suggs, err = listMealsByScore(q)
	if err != nil {
		return nil, err
	}

	return suggs, nil
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
*/
