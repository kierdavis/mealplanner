package mpdb

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"sort"
	"time"
)

const CalculateTagScoresSQL = "SELECT tag.tag, ABS(DATEDIFF(serving.dateserved, ?)) " +
	"FROM tag " +
	"INNER JOIN serving ON serving.mealid = tag.mealid"

const ListSuggestionsSQL = "SELECT meal.id, meal.name, meal.recipe, meal.favourite, MIN(ABS(DATEDIFF(serving.dateserved, ?))) " +
	"FROM meal " +
	"LEFT JOIN serving ON meal.id = serving.mealid " +
	"GROUP BY meal.id"

func GenerateSuggestions(q Queryable, date time.Time) (suggs []*mpdata.Suggestion, err error) {
	s := mpdata.NewScorer()

	err = calculateTagScores(q, s, date)
	if err != nil {
		return nil, err
	}

	suggs, err = listSuggestions(q, date)
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

	return suggs, nil
}

func calculateTagScores(q Queryable, s *mpdata.Scorer, date time.Time) (err error) {
	rows, err := q.Query(CalculateTagScoresSQL, date)
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

func listSuggestions(q Queryable, date time.Time) (suggs []*mpdata.Suggestion, err error) {
	rows, err := q.Query(ListSuggestionsSQL, date)
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

		if csd.Valid {
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

func getTagsForSuggestions(q Queryable, suggs []*mpdata.Suggestion) (err error) {
	getTagsStmt, err := q.Prepare(GetMealTagsSQL)
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
