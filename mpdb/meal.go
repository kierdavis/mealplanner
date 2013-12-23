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

const GetMealSQL =
	"SELECT meal.name, meal.recipe, meal.favourite " +
	"FROM meal " +
	"WHERE meal.id = ?"

const GetMealTagsSQL =
	"SELECT tag.tag " +
	"FROM tag " +
	"WHERE tag.mealid = ?"

const AddMealSQL =
	"INSERT INTO meal " +
	"VALUES (NULL, ?, ?, ?)"

const UpdateMealSQL =
	"UPDATE meal " +
	"SET meal.name = ? " +
	"    meal.recipe = ? " +
	"    meal.favourite = ? " +
	"WHERE meal.id = ?"

const DeleteAllMealTagsSQL =
	"DELETE FROM tag " +
	"WHERE tag.mealid = ?"

const AddMealTagsSQL =
	"INSERT INTO tag " +
	"VALUES (?, ?)"

// ListMeals fetches and returns a list of all meals in the database. If its
// parameter 'sortByName' is true, the meals are sorted in alphabetical order
// by name.
func ListMeals(q Queryable, sortByName bool) (meals []*mpdata.Meal, err error) {
	var query string
	if sortByName {
		query = ListMealsByNameSQL
	} else {
		query = ListMealsSQL
	}
	
	rows, err := q.Query(query)
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

func GetMeal(q Queryable, mealID uint64) (meal *mpdata.Meal, err error) {
	meal = &mpdata.Meal{ID: mealID}
	err = q.QueryRow(GetMealSQL, mealID).Scan(&meal.Name, &meal.RecipeURL, &meal.Favourite)
	if err != nil {
		return nil, err
	}
	
	return meal, nil
}

func GetMealTags(q Queryable, mealID uint64) (tags []string, err error) {
	rows, err := q.Query(GetMealTagsSQL, mealID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tag string
	
	for rows.Next() {
		err = rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		
		tags = append(tags, tag)
	}
	
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	
	return tags, nil
}

func GetMealWithTags(q Queryable, mealID uint64) (meal *mpdata.Meal, err error) {
	meal, err = GetMeal(q, mealID)
	if err != nil {
		return nil, rer
	}
	
	meal.Tags, err = GetMealTags(q, mealID)
	if err != nil {
		return nil, err
	}
	
	meal.HasTags = true
	return meal, nil
}

func AddMeal(q Queryable, meal *mpdata.Meal) (mealID uint64, err error) {
	result, err := q.Exec(AddMealSQL, meal.Name, meal.RecipeURL, meal.Favourite)
	if err != nil {
		return 0, err
	}
	
	n, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	
	return uint64(n), nil
}

func AddMealTags(q Queryable, mealID uint64, tags []string) (err error) {
	stmt, err := q.Prepare(AddMealTagsSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	for _, tag := range tags {
		_, err = stmt.Exec(mealID, tag)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func AddMealWithTags(q Queryable, meal *mpdata.Meal) (mealID uint64, err error) {
	if !meal.HasTags {
		return 0, fmt.Errorf("meal argument does not have an attached tag list")
	}
	
	mealID, err = AddMeal(q, meal)
	if err != nil {
		return 0, err
	}
	
	err = AddMealTags(q, meal.ID, meal.Tags)
	if err != nil {
		return 0, err
	}
	
	return mealID
}

func UpdateMeal(q Queryable, meal *mpdata.Meal) (err error) {
	_, err = q.Exec(UpdateMealSQL, meal.Name, meal.RecipeURL, meal.Favourite, meal.ID)
	return err
}

func DeleteAllMealTags(q Queryable, mealID uint64) (err error) {
	_, err = q.Exec(DeleteAllMealTagsSQL, mealID)
	return err
}

func UpdateMealTags(q Queryable, mealID uint64, tags []string) (err error) {
	err = DeleteAllMealTags(q, mealID)
	if err != nil {
		return err
	}
	
	return AddMealTags(q, mealID, tags)
}

func UpdateMealWithTags(q Queryable, meal *mpdata.Meal) (err error) {
	if !meal.HasTags {
		return 0, fmt.Errorf("meal argument does not have an attached tag list")
	}
	
	err = UpdateMeal(q, meal)
	if err != nil {
		return err
	}
	
	return UpdateMealTags(q, meal.ID, meal.Tags)
}
