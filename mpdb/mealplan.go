package mpdb

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"time"
)

// SQL statement for retrieving information about a meal plan.
const GetMealPlanSQL = "SELECT mealplan.notes, mealplan.startdate, mealplan.enddate FROM mealplan WHERE mealplan.id = ?"

// SQL statement for adding a meal plan to the database.
const AddMealPlanSQL = "INSERT INTO mealplan VALUES (NULL, ?, ?, ?)"

// SQL statement for setting the notes associated with a meal plan.
const UpdateNotesSQL = "UPDATE mealplan SET mealplan.notes = ? WHERE mealplan.id = ?"

// SQL statement for deleting a meal plan identified by its ID.
const DeleteMealPlanSQL = "DELETE FROM mealplan WHERE mealplan.id = ?"

// SQL statement for listing meal plans that overlap a given date range.
const ListMealPlansBetweenSQL = "SELECT mealplan.id, mealplan.startdate, mealplan.enddate " +
	"FROM mealplan " +
	"WHERE ? <= mealplan.enddate && mealplan.startdate <= ?"

// SQL statement for retrieving information about a meal serving.
const GetServingSQL = "SELECT serving.mealid FROM serving WHERE serving.mealplanid = ? AND serving.dateserved = ?"

// SQL statement for retrieving the servings associated with a meal plan.
const GetServingsSQL = "SELECT serving.dateserved, serving.mealid FROM serving WHERE serving.mealplanid = ?"

// SQL statement for returning the number of servings on a meal plan.
const CountServingsSQL = "SELECT COUNT(serving.dateserved) FROM serving WHERE serving.mealplanid = ?"

// SQL statement for deleting a serving.
const DeleteServingSQL = "DELETE FROM serving WHERE serving.mealplanid = ? AND serving.dateserved = ?"

// SQL statement for deleting all servings in a meal plan.
const DeleteServingsSQL = "DELETE FROM serving WHERE serving.mealplanid = ?"

// SQL statement for deleting all servings of a given meal.
const DeleteServingsOfSQL = "DELETE FROM serving WHERE serving.mealid = ?"

// SQL statement for adding a meal serving.
const InsertServingSQL = "INSERT INTO serving VALUES (?, ?, ?)"

// GetMealPlan returns information about the meal plan identified by 'mpID'.
func GetMealPlan(q Queryable, mpID uint64) (mp *mpdata.MealPlan, err error) {
	mp = &mpdata.MealPlan{ID: mpID}
	err = q.QueryRow(GetMealPlanSQL, mpID).Scan(&mp.Notes, &mp.StartDate, &mp.EndDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return mp, nil
}

// AddMealPlan adds the information contained in 'mp' to the database as a new
// meal plan record. It assigns the identifier of the newly created record to
// the ID field of the meal plan.
func AddMealPlan(q Queryable, mp *mpdata.MealPlan) (err error) {
	result, err := q.Exec(AddMealPlanSQL, mp.Notes, mp.StartDate, mp.EndDate)
	if err != nil {
		return err
	}

	mpID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	mp.ID = uint64(mpID)

	return nil
}

// UpdateNotes sets the notes associated with the meal plan identified by 'mpID'
// to 'notes'.
func UpdateNotes(q Queryable, mpID uint64, notes string) (err error) {
	_, err = q.Exec(UpdateNotesSQL, notes, mpID)
	return err
}

// DeleteMealPlan deletes the meal plan record identified by 'mpID'. If no such
// meal plan exists, no error is raised.
func DeleteMealPlan(q Queryable, mpID uint64) (err error) {
	_, err = q.Exec(DeleteMealPlanSQL, mpID)
	return err
}

// ListMealPlansBetween returns a list of all meal plans in the database whose
// date range (start date to end date) overlaps with the given date range
// ('from' to 'to').
func ListMealPlansBetween(q Queryable, from time.Time, to time.Time) (mps []*mpdata.MealPlan, err error) {
	rows, err := q.Query(ListMealPlansBetweenSQL, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		mp := &mpdata.MealPlan{}

		err = rows.Scan(&mp.ID, &mp.StartDate, &mp.EndDate)
		if err != nil {
			return nil, err
		}

		mps = append(mps, mp)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return mps, nil
}

// GetServing returns information about the meal serving identified by the
// meal plan identifier 'mpID' and the serving date 'date'.
func GetServing(q Queryable, mpID uint64, date time.Time) (serving *mpdata.Serving, err error) {
	serving = &mpdata.Serving{MealPlanID: mpID, Date: date}
	err = q.QueryRow(GetServingSQL, mpID, date).Scan(&serving.MealID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return serving, nil
}

// GetServings returns a slice containing the servings that are part of the
// meal plan identified by 'mpID'.
func GetServings(q Queryable, mpID uint64) (servings []*mpdata.Serving, err error) {
	rows, err := q.Query(GetServingsSQL, mpID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		serving := &mpdata.Serving{MealPlanID: mpID}

		err = rows.Scan(&serving.Date, &serving.MealID)
		if err != nil {
			return nil, err
		}

		servings = append(servings, serving)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return servings, nil
}

// CountServings returns the number of servings in the meal plan identified by
// 'mpID'.
func CountServings(q Queryable, mpID uint64) (numServings int, err error) {
	err = q.QueryRow(CountServingsSQL, mpID).Scan(&numServings)
	if err != nil {
		return 0, err
	}

	return numServings, nil
}

// GetMealPlanWithServings returns the information about the meal plan
// identified by 'mpID' including its servings.
func GetMealPlanWithServings(q Queryable, mpID uint64) (mps *mpdata.MealPlanWithServings, err error) {
	mp, err := GetMealPlan(q, mpID)
	if err != nil {
		return nil, err
	}

	servings, err := GetServings(q, mpID)
	if err != nil {
		return nil, err
	}

	mps = &mpdata.MealPlanWithServings{
		MealPlan: mp,
		Servings: servings,
	}
	return mps, nil
}

// DeleteServing deletes the serving at 'date' on the meal plan identified by
// 'mpID'. If no such serving exists, no error is raised.
func DeleteServing(q Queryable, mpID uint64, date time.Time) (err error) {
	_, err = q.Exec(DeleteServingSQL, mpID, date)
	return err
}

// DeleteServings deletes all servings on the meal plan identified by 'mpID'. If
// no such servings exist, no error is raised.
func DeleteServings(q Queryable, mpID uint64) (err error) {
	_, err = q.Exec(DeleteServingsSQL, mpID)
	return err
}

// DeleteServingsOf deletes all servings of the meal identified by 'mealID'. IF
// no such servings exist, no error is raised.
func DeleteServingsOf(q Queryable, mealID uint64) (err error) {
	_, err = q.Exec(DeleteServingsOfSQL, mealID)
	return err
}

// AddServing adds the information containing in 'serving' to a new serving
// record in the database.
func AddServing(q Queryable, serving *mpdata.Serving) (err error) {
	_, err = q.Exec(InsertServingSQL, serving.MealPlanID, serving.Date, serving.MealID)
	return err
}

// AutoFillMealPlan assigns servings to every day in 'mp' using the top
// suggestion for each day.
func AutoFillMealPlan(q Queryable, mp *mpdata.MealPlan) (err error) {
	for _, date := range mp.Days() {
		err = AutoFillMealPlanDay(q, mp.ID, date)
		if err != nil {
			return err
		}
	}
	
	return nil
}

// AutoFillMealPlanDay assigns a serving to day 'date' on the meal plan
// identified by 'mpID' using the top suggestion.
func AutoFillMealPlanDay(q Queryable, mpID uint64, date time.Time) (err error) {
	suggs, err := GenerateSuggestions(q, date)
	if err != nil {
		return err
	}
	
	err = DeleteServing(q, mpID, date)
	if err != nil {
		return err
	}
	
	serving := &mpdata.Serving{
		MealPlanID: mpID,
		Date: date,
		MealID: suggs[0].MT.Meal.ID,
	}
	
	return AddServing(q, serving)
}
