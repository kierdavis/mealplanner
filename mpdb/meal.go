package mpdb

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
)

// SQL statement for listing meals.
const ListMealsSQL = "SELECT meal.id, meal.name, meal.recipe, meal.favourite FROM meal"

// SQL statement for listing meals sorted by name.
const ListMealsByNameSQL = "SELECT meal.id, meal.name, meal.recipe, meal.favourite FROM meal ORDER BY meal.name ASC"

// SQL statement for fetching information about a meal.
const GetMealSQL = "SELECT meal.name, meal.recipe, meal.favourite FROM meal WHERE meal.id = ?"

// SQL statement for fetching tags associated with a meal.
const GetMealTagsSQL = "SELECT tag.tag FROM tag WHERE tag.mealid = ?"

// SQL statement for adding a meal.
const AddMealSQL = "INSERT INTO meal VALUES (NULL, ?, ?, ?)"

// SQL statement for updating the information about a meal.
const UpdateMealSQL = "UPDATE meal SET meal.name = ?, meal.recipe = ?, meal.favourite = ? WHERE meal.id = ?"

// SQL statement for deleting all tags associated with a meal.
const DeleteMealTagsSQL = "DELETE FROM tag WHERE tag.mealid = ?"

// SQL statement for adding a tag to a meal.
const AddMealTagSQL = "INSERT INTO tag VALUES (?, ?)"

// SQL statement for testing whether a meal is marked as a favourite.
const IsFavouriteSQL = "SELECT meal.favourite FROM meal WHERE meal.id = ?"

// SQL statement to set the "favourite" status of a meal.
const SetFavouriteSQL = "UPDATE meal SET meal.favourite = ? WHERE meal.id = ?"

// SQL statement to delete a meal.
const DeleteMealSQL = "DELETE FROM meal WHERE meal.id = ?"

// SQL statement to list all tags in the database.
const ListAllTagsSQL = "SELECT DISTINCT tag.tag FROM tag"

// SQL statement to list all tags in the database sorted by name.
const ListAllTagsByNameSQL = "SELECT DISTINCT tag.tag FROM tag ORDER BY tag.tag ASC"

// ListMeals fetches and returns a list of all meals in the database. If the
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

// ListMealsWithTags fetches and returns a list of all meals in the database
// with their associated tags. If the parameter 'sortByName' is true, the meals
// are sorted in alphabetical order by name.
func ListMealsWithTags(q Queryable, sortByName bool) (mts []mpdata.MealWithTags, err error) {
	var query string
	if sortByName {
		query = ListMealsByNameSQL
	} else {
		query = ListMealsSQL
	}

	getTagsStmt, err := q.Prepare(GetMealTagsSQL)
	if err != nil {
		return nil, err
	}
	defer getTagsStmt.Close()

	rows, err := q.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		mt := mpdata.MealWithTags{
			Meal: &mpdata.Meal{},
		}

		err = rows.Scan(&mt.Meal.ID, &mt.Meal.Name, &mt.Meal.RecipeURL, &mt.Meal.Favourite)
		if err != nil {
			return nil, err
		}
		
		mt.Tags, err = getMealTagsPrepared(q, getTagsStmt, mt.Meal.ID)
		if err != nil {
			return nil, err
		}

		mts = append(mts, mt)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return mts, nil
}

// GetMeal fetches information from the database about the meal identified by
// 'mealID'.
func GetMeal(q Queryable, mealID uint64) (meal *mpdata.Meal, err error) {
	meal = &mpdata.Meal{ID: mealID}
	err = q.QueryRow(GetMealSQL, mealID).Scan(&meal.Name, &meal.RecipeURL, &meal.Favourite)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		
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

	return readTags(rows)
}

// getMealTagsPrepared fetches the list of tags associated with the meal
// identified by 'mealID'. 'stmt' is assumed to be a prepared statement compiled
// from GetMealTagsSQL.
func getMealTagsPrepared(q Queryable, stmt *sql.Stmt, mealID uint64) (tags []string, err error) {
	rows, err := stmt.Query(mealID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return readTags(rows)
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

// DeleteMealTags deletes all tags in the database associated with the meal
// identified by 'mealID'.
func DeleteMealTags(q Queryable, mealID uint64) (err error) {
	_, err = q.Exec(DeleteMealTagsSQL, mealID)
	return err
}

// UpdateMealTags replaces the tags associated with the meal identified by
// 'mealID' with the list given by 'tags'.
func UpdateMealTags(q Queryable, mealID uint64, tags []string) (err error) {
	err = DeleteMealTags(q, mealID)
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

// ToggleFavourite toggles the "favourite" status of the meal identified by
// 'mealID', and returns the new favourite status.
func ToggleFavourite(q Queryable, mealID uint64) (isFavourite bool, err error) {
	err = q.QueryRow(IsFavouriteSQL, mealID).Scan(&isFavourite)
	if err != nil {
		return false, err
	}
	
	isFavourite = !isFavourite
	_, err = q.Exec(SetFavouriteSQL, isFavourite, mealID)
	return isFavourite, err
}

// DeleteMeal deletes the meal record identified by 'mealID'.
func DeleteMeal(q Queryable, mealID uint64) (err error) {
	_, err = q.Exec(DeleteMealSQL, mealID)
	return err
}

// DeleteMealWithTags deletes the meal record identified by 'mealID', and all
// tag records associated with it.
func DeleteMealWithTags(q Queryable, mealID uint64) (err error) {
	err = DeleteMeal(q, mealID)
	if err != nil {
		return err
	}
	
	return DeleteMealTags(q, mealID)
}

// ListAllTags returns a list (without duplicates) of all tags that appear in
// the database. If the 'sortByName' parameter is true, the tags are sorted into
// alphabetical order.
func ListAllTags(q Queryable, sortByName bool) (tags []string, err error) {
	var query string
	if sortByName {
		query = ListAllTagsByNameSQL
	} else {
		query = ListAllTagsSQL
	}
	
	rows, err := q.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return readTags(rows)
}

// readTags reads a *sql.Rows as produced by GetMealTags or
// getMealTagsPrepared and converts it into a slice of tags.
func readTags(rows *sql.Rows) (tags []string, err error) {
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
