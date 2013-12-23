package mpdb

import (
	"github.com/kierdavis/mealplanner/mpdata"
)

// SQL statement for listing meals.
const ListMealsSQL = "SELECT meal.id, meal.name, meal.recipe, meal.favourite " +
	"FROM meal"

// SQL statement for listing meals sorted by name.
const ListMealsByNameSQL = "SELECT meal.id, meal.name, meal.recipe, meal.favourite " +
	"FROM meal " +
	"ORDER BY meal.name ASC"

// SQL statement for fetching information about a meal.
const GetMealSQL = "SELECT meal.name, meal.recipe, meal.favourite " +
	"FROM meal " +
	"WHERE meal.id = ?"

// SQL statement for fetching tags associated with a meal.
const GetMealTagsSQL = "SELECT tag.tag " +
	"FROM tag " +
	"WHERE tag.mealid = ?"

// SQL statement for adding a meal.
const AddMealSQL = "INSERT INTO meal " +
	"VALUES (NULL, ?, ?, ?)"

// SQL statement for updating the information about a meal.
const UpdateMealSQL = "UPDATE meal " +
	"SET meal.name = ? " +
	"    meal.recipe = ? " +
	"    meal.favourite = ? " +
	"WHERE meal.id = ?"

// SQL statement for deletting all tags associated with a meal.
const DeleteAllMealTagsSQL = "DELETE FROM tag " +
	"WHERE tag.mealid = ?"

// SQL statement for adding a tag to a meal.
const AddMealTagSQL = "INSERT INTO taS" +
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

// GetMeal fetches information from the database about the meal identified by
// 'mealID'.
func GetMeal(q Queryable, mealID uint64) (meal *mpdata.Meal, err error) {
	meal = &mpdata.Meal{ID: mealID}
	err = q.QueryRow(GetMealSQL, mealID).Scan(&meal.Name, &meal.RecipeURL, &meal.Favourite)
	if err != nil {
		return nil, err
	}

	return meal, nil
}

// GetMealTags fetches the list of tags associated with the meal identified by
// 'mealID'.
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

// GetMealWithTags combines GetMeal and GetMealTags.
func GetMealWithTags(q Queryable, mealID uint64) (mt mpdata.MealWithTags, err error) {
	meal, err := GetMeal(q, mealID)
	if err != nil {
		return mpdata.MealWithTags{}, err
	}

	tags, err := GetMealTags(q, mealID)
	if err != nil {
		return mpdata.MealWithTags{}, err
	}

	return mpdata.MealWithTags{meal, tags}, nil
}

// AddMeal adds the information in 'meal' to the database as a new record, then
// sets 'meal.ID' to the identifier of this new record.
func AddMeal(q Queryable, meal *mpdata.Meal) (err error) {
	result, err := q.Exec(AddMealSQL, meal.Name, meal.RecipeURL, meal.Favourite)
	if err != nil {
		return err
	}

	mealID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	meal.ID = uint64(mealID)

	return nil
}

// AddMealTags adds the the list of tags given in 'tags' to the meal identified
// by 'mealID'.
func AddMealTags(q Queryable, mealID uint64, tags []string) (err error) {
	stmt, err := q.Prepare(AddMealTagSQL)
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

// AddMealWithTags combines 'AddMeal' and 'AddMealTags'.
func AddMealWithTags(q Queryable, mt mpdata.MealWithTags) (err error) {
	err = AddMeal(q, mt.Meal)
	if err != nil {
		return err
	}

	return AddMealTags(q, mt.Meal.ID, mt.Tags)
}

// UpdateMeal replaces with the information in the database for the meal
// identified by 'meal.ID' with the information in 'meal'.
func UpdateMeal(q Queryable, meal *mpdata.Meal) (err error) {
	_, err = q.Exec(UpdateMealSQL, meal.Name, meal.RecipeURL, meal.Favourite, meal.ID)
	return err
}

// DeleteAllMealTags deletes all tags in the database associated with the meal
// identified by 'mealID'.
func DeleteAllMealTags(q Queryable, mealID uint64) (err error) {
	_, err = q.Exec(DeleteAllMealTagsSQL, mealID)
	return err
}

// UpdateMealTags replaces the tags associated with the meal identified by
// 'mealID' with the list given by 'tags'.
func UpdateMealTags(q Queryable, mealID uint64, tags []string) (err error) {
	err = DeleteAllMealTags(q, mealID)
	if err != nil {
		return err
	}

	return AddMealTags(q, mealID, tags)
}

// UpdateMealWithTags combines UpdateMeal and UpdateMealTags.
func UpdateMealWithTags(q Queryable, mt mpdata.MealWithTags) (err error) {
	err = UpdateMeal(q, mt.Meal)
	if err != nil {
		return err
	}

	return UpdateMealTags(q, mt.Meal.ID, mt.Tags)
}
