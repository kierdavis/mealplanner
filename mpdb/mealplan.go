package mpdb

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
	"time"
)

const GetMealPlanSQL = "SELECT mealplan.notes, mealplan.startdate, mealplan.enddate FROM mealplan WHERE mealplan.id = ?"

const AddMealPlanSQL = "INSERT INTO mealplan VALUES (NULL, ?, ?, ?)"

const GetServingSQL = "SELECT serving.mealid FROM serving WHERE serving.mealplanid = ? AND serving.dateserved = ?"

const GetServingsSQL = "SELECT serving.dateserved, serving.mealid FROM serving WHERE serving.mealplanid = ?"

const DeleteServingSQL = "DELETE FROM serving WHERE serving.mealid = ? AND serving.dateserved = ?"

const InsertServingSQL = "INSERT INTO serving VALUES (?, ?, ?)"

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

func DeleteServing(q Queryable, mpID uint64, date time.Time) (err error) {
	_, err = q.Exec(DeleteServingSQL, mpID, date)
	return err
}

func AddServing(q Queryable, serving *mpdata.Serving) (err error) {
	_, err = q.Exec(InsertServingSQL, serving.MealPlanID, serving.Date, serving.MealID)
	return err
}
