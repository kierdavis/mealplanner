package mpdb

import (
	"github.com/kierdavis/mealplanner/mpdata"
)

const ListMealsSQL =
	"SELECT meal.id, meal.name, meal.recipe, meal.favourite " +
	"FROM meal"

const ListMealsByNameSQL =
	"SELECT meal.id, meal.name, meal.recipe, meal.favourite " +
	"FROM meal " +
	"ORDER BY meal.name ASC"

// ListMeals fetches and returns a list of all meals in the database. If its
// parameter 'sortByName' is true, the meals are sorted in alphabetical order
// by name.
func (db DB) ListMeals(sortByName bool) (meals []*mpdata.Meal, err error) {
	var query string
	if sortByName {
		query = ListMealsByNameSQL
	} else {
		query = ListMealsSQL
	}
	
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		meal := &mpdata.Meal{}
		
		err = rows.Scan(&meal.ID, &meal.Name, &meal.RecipeURL, &meal.Favourite)
		if err != nil {
			return nil, err
		}
		
		meals = append(meals, meal)
	}
	
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	
	return meals, nil
}
