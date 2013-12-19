package mpdb

// A pair consisting of a meal identifier and a corresponding score
type MealScore struct {
	MealID uint64
	Score float32
}

const SuggestionBatchSize = 10

func (db DB) GenerateSuggestions(date time.Time) (suggs []*mpdata.Meal, err error) {
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
	csdStmt, err := db.prepareCsdStmt()
	if err != nil {
		return nil, err
	}
	defer csdStmt.Close() // Defer cleanup of the prepared statement
	
	// Prepare findNearbyTagDists query for repeated use
	ntdStmt, err := db.prepareNtdStmt()
	if err != nil {
		return nil, err
	}
	defer ntdStmt.Close() // Defer cleanup of the prepared statement
	
	// Create buffer to hold current batch in
	scorePairs := make([]MealScore, 0, SuggestionBatchSize)
	
	for _, meal := range meals {
		// Find closest serving distance
		dist, err := db.findClosestServingDistance(csdStmt, meal.ID, date)
		if err != nil {
			return nil, err
		}
		
		// Find nearby tag distances
		tds, err := db.findNearbyTagDists(ntdStmt, date)
		if err != nil {
			return nil, err
		}
		
		// Calculate score and add to current batch
		score := mpdata.CalculateScore(meal.Favourite, dist, tds)
		scorePair := MealScore{meal.ID, score}
		scorePairs = append(scorePairs, scorePair)
		
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
	}
	
	// If there are scores not yet inserted,
	if len(scorePairs) > 0 {
		// Insert them
		err = db.insertScores(scorePairs)
		if err != nil {
			return nil, err
		}
	}
	
	// List all meals, but sorted by score
	suggs, err = db.listMealsByScore()
	if err != nil {
		return nil, err
	}
	
	return suggs, nil
}

const CreateScoreTableSQL =
	"CREATE TEMPORARY TABLE score (" +
	"  mealid BIGINT UNSIGNED NOT NULL, " +
	"  score FLOAT NOT NULL" +
	")"

func (db DB) createScoreTable() (err error) {
	_, err = db.Exec(CreateScoreTableSQL)
	return err
}

const DropScoreTableSQL =
	"DROP TABLE score"

func (db DB) dropScoreTable() (err error) {
	_, err = db.Exec(DropScoreTableSQL)
	return err
}
