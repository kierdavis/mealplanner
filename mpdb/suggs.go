package mpdb

// A pair consisting of a meal identifier and a corresponding score
type MealScore struct {
	MealID uint64
	Score float32
}

const SCORE_BUFFER_SIZE = 32

func (db DB) GenerateSuggestions(date time.Time) (suggs []*mpdata.Meal, err error) {
	err = db.createScoreTable()
	if err != nil {
		return nil, err
	}
	defer db.dropScoreTable()
	
	meals, err := db.ListMeals(false)
	if err != nil {
		return nil, err
	}
	
	csdStmt, err := db.prepareCsdStmt()
	if err != nil {
		return nil, err
	}
	defer csdStmt.Close()
	
	ntdStmt, err := db.prepareNtdStmt()
	if err != nil {
		return nil, err
	}
	defer ntdStmt.Close()
	
	scorePairs := make([]MealScore, 0, SCORE_BUFFER_SIZE)
	
	for _, meal := range meals {
		dist, err := db.findClosestServingDistance(csdStmt, meal.ID, date)
		if err != nil {
			return nil, err
		}
		
		tds, err := db.findNearbyTagDists(ntdStmt, date)
		if err != nil {
			return nil, err
		}
		
		score := mpdata.CalculateScore(meal.Favourite, dist, tds)
		scorePair := MealScore{meal.ID, score}
		scorePairs = append(scorePairs, scorePair)
		
		if len(scorePairs) == cap(scorePairs) {
			err = db.insertScores(scorePairs)
			if err != nil {
				return nil, err
			}
			
			scorePairs = scorePairs[:0]
		}
	}
	
	if len(scorePairs) > 0 {
		err = db.insertScores(scorePairs)
		if err != nil {
			return nil, err
		}
	}
	
	suggs, err = db.listMealsByScore()
	if err != nil {
		return nil, err
	}
	
	return suggs, nil
}
