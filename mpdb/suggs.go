package mpdb

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"sort"
	"time"
)

// GenerateSuggestions calculates a score for each meal in the database based on
// their suitability for serving on 'date'. These are returned as a list of
// Suggestions.
func GenerateSuggestions(q Queryable, mpID uint64, date time.Time) (suggs []*mpdata.Suggestion, err error) {
	s := mpdata.NewScorer()

	err = calculateTagScores(q, s, mpID, date)
	if err != nil {
		return nil, err
	}

	suggs, err = listSuggestions(q, mpID, date)
	if err != nil {
		return nil, err
	}

	err = getTagsForSuggestions(q, suggs)
	if err != nil {
		return nil, err
	}

	for _, sugg := range suggs {
		s.ScoreSuggestion(sugg)
	}

	sort.Sort(mpdata.SuggestionSlice(suggs))

	// Scale the scores to between 0 and 1
	maxScore := suggs[0].Score
	minScore := suggs[len(suggs)-1].Score
	scoreRange := maxScore - minScore

	for _, sugg := range suggs {
		sugg.Score = (sugg.Score - minScore) / scoreRange
	}

	return suggs, nil
}

// calculateTagScores prepares the scorer 's' by adding a score for each usage
// of a tag.
func calculateTagScores(q Queryable, s *mpdata.Scorer, mpID uint64, date time.Time) (err error) {
	rows, err := q.Query("SELECT tag.tag, MIN(ABS(DATEDIFF(serving.dateserved, ?))) FROM tag INNER JOIN serving ON serving.mealid = tag.mealid GROUP BY tag.tag", date)
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

// listSuggestions returns a list of meals (without tags) and the distance
// between 'date' and their closest serving to 'date'.
func listSuggestions(q Queryable, mpID uint64, date time.Time) (suggs []*mpdata.Suggestion, err error) {
	rows, err := q.Query("SELECT meal.id, meal.name, meal.recipe, meal.favourite, MIN(ABS(DATEDIFF(serving.dateserved, ?))) FROM meal LEFT JOIN serving ON meal.id = serving.mealid GROUP BY meal.id", date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		meal := new(mpdata.Meal)
		sugg := new(mpdata.Suggestion)
		sugg.MT.Meal = meal

		var csd sql.NullInt64

		err = rows.Scan(&meal.ID, &meal.Name, &meal.RecipeURL, &meal.Favourite, &csd)
		if err != nil {
			return nil, err
		}

		if csd.Valid && csd.Int64 != 0 {
			sugg.CSD = int(csd.Int64)
		} else {
			sugg.CSD = -1
		}

		suggs = append(suggs, sugg)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return suggs, nil
}

// getTagsForSuggestions fills the tags field of each suggestion in 'suggs'.
func getTagsForSuggestions(q Queryable, suggs []*mpdata.Suggestion) (err error) {
	getTagsStmt, err := q.Prepare("SELECT tag.tag FROM tag WHERE tag.mealid = ?")
	if err != nil {
		return err
	}
	defer getTagsStmt.Close()

	for _, sugg := range suggs {
		sugg.MT.Tags, err = getMealTagsPrepared(getTagsStmt, sugg.MT.Meal.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
