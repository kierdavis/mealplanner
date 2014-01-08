package mpdb

import (
	"github.com/kierdavis/mealplanner/mpdata"
)

const GetMealPlanSQL = "SELECT mealplan.notes, mealplan.startdate, mealplan.enddate FROM mealplan WHERE mealplan.id = ?"

const AddMealPlanSQL = "INSERT INTO mealplan VALUES (NULL, ?, ?, ?)"

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
