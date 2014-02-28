package mpdb

import (
	"database/sql"
	"github.com/kierdavis/mealplanner/mpdata"
)

const SearchTextFunc = "CONCAT(meal.name, ' ', meal.recipe, ' ', IFNULL((SELECT GROUP_CONCAT(tag.tag SEPARATOR ' ') FROM tag WHERE tag.mealid = meal.id), ''))"

// SQL statement for listing meals.
const ListMealsSQL = "SELECT meal.id, meal.name, meal.recipe, meal.favourite FROM meal"

// SQL statement for listing meals sorted by name.
const ListMealsByNameSQL = ListMealsSQL + " ORDER BY meal.name"

const CreateSearchPatternsTableSQL = "CREATE TEMPORARY TABLE search_patterns (pattern VARCHAR(255))"

const InsertSearchPatternSQL = "INSERT INTO search_patterns VALUES (?)"

const DropSearchPatternsTableSQL = "DROP TABLE search_patterns"

// SQL statement for fetching information about a meal.
const GetMealSQL = "SELECT meal.name, meal.recipe, meal.favourite FROM meal WHERE meal.id = ?"

// SQL statement for fetching tags associated with a meal.
const GetMealTagsSQL = "SELECT tag.tag FROM tag WHERE tag.mealid = ?"

const UpdateSearchTextSQL = "UPDATE meal SET meal.searchtext = " + SearchTextFunc + " WHERE meal.id = ?"

// SQL statement for adding a meal.
const AddMealSQL = "INSERT INTO meal VALUES (NULL, ?, ?, ?, '')"

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

	return readMeals(rows)
}

// ListMealsWithTags fetches and returns a list of all meals in the database
// with their associated tags. If the parameter 'sortByName' is true, the meals
// are sorted in alphabetical order by name.
func ListMealsWithTags(q Queryable, sortByName bool) (mts []mpdata.MealWithTags, err error) {
	meals, err := ListMeals(q, sortByName)
	if err != nil {
		return nil, err
	}

	return AttachMealTags(q, meals)
}

func SearchMeals(q Queryable, words []string, sortByName bool) (meals []*mpdata.Meal, err error) {
	query := "SELECT meal.id, meal.name, meal.recipe, meal.favourite FROM meal"
	conjuctive := "WHERE"
	args := make([]interface{}, len(words))
	for i, word := range words {
		query += " " + conjuctive + " meal.searchtext LIKE ?"
		args[i] = "%" + word + "%"
		conjuctive = "AND"
	}

	if sortByName {
		query += " ORDER BY meal.name"
	}

	rows, err := q.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return readMeals(rows)
}

func SearchMealsWithTags(q Queryable, words []string, sortByName bool) (mts []mpdata.MealWithTags, err error) {
	meals, err := SearchMeals(q, words, sortByName)
	if err != nil {
		return nil, err
	}

	return AttachMealTags(q, meals)
}

func readMeals(rows *sql.Rows) (meals []*mpdata.Meal, err error) {
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

func AttachMealTags(q Queryable, meals []*mpdata.Meal) (mts []mpdata.MealWithTags, err error) {
	getTagsStmt, err := q.Prepare(GetMealTagsSQL)
	if err != nil {
		return nil, err
	}
	defer getTagsStmt.Close()

	for _, meal := range meals {
		tags, err := getMealTagsPrepared(getTagsStmt, meal.ID)
		if err != nil {
			return nil, err
		}

		mt := mpdata.MealWithTags{
			Meal: meal,
			Tags: tags,
		}
		mts = append(mts, mt)
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
func getMealTagsPrepared(stmt *sql.Stmt, mealID uint64) (tags []string, err error) {
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

	return mpdata.MealWithTags{Meal: meal, Tags: tags}, nil
}

// AddMeal adds the information in 'meal' to the database as a new record, then
// sets 'meal.ID' to the identifier of this new record.
func AddMeal(q Queryable, meal *mpdata.Meal) (err error) {
	err = addMeal(q, meal)
	if err != nil {
		return err
	}
	return UpdateSearchText(q, meal.ID)
}

func addMeal(q Queryable, meal *mpdata.Meal) (err error) {
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

func UpdateSearchText(q Queryable, mealID uint64) (err error) {
	_, err = q.Exec(UpdateSearchTextSQL, mealID)
	return err
}

// AddMealTags adds the the list of tags given in 'tags' to the meal identified
// by 'mealID'.
func AddMealTags(q Queryable, mealID uint64, tags []string) (err error) {
	err = addMealTags(q, mealID, tags)
	if err != nil {
		return err
	}
	return UpdateSearchText(q, mealID)
}

func addMealTags(q Queryable, mealID uint64, tags []string) (err error) {
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
	err = addMealWithTags(q, mt)
	if err != nil {
		return err
	}
	return UpdateSearchText(q, mt.Meal.ID)
}

func addMealWithTags(q Queryable, mt mpdata.MealWithTags) (err error) {
	err = addMeal(q, mt.Meal)
	if err != nil {
		return err
	}

	return addMealTags(q, mt.Meal.ID, mt.Tags)
}

// UpdateMeal replaces with the information in the database for the meal
// identified by 'meal.ID' with the information in 'meal'.
func UpdateMeal(q Queryable, meal *mpdata.Meal) (err error) {
	err = updateMeal(q, meal)
	if err != nil {
		return err
	}
	return UpdateSearchText(q, meal.ID)
}

func updateMeal(q Queryable, meal *mpdata.Meal) (err error) {
	_, err = q.Exec(UpdateMealSQL, meal.Name, meal.RecipeURL, meal.Favourite, meal.ID)
	return err
}

// DeleteMealTags deletes all tags in the database associated with the meal
// identified by 'mealID'. If no such tags exist, no error is raised.
func DeleteMealTags(q Queryable, mealID uint64) (err error) {
	err = deleteMealTags(q, mealID)
	if err != nil {
		return err
	}
	return UpdateSearchText(q, mealID)
}

func deleteMealTags(q Queryable, mealID uint64) (err error) {
	_, err = q.Exec(DeleteMealTagsSQL, mealID)
	return err
}

// UpdateMealTags replaces the tags associated with the meal identified by
// 'mealID' with the list given by 'tags'.
func UpdateMealTags(q Queryable, mealID uint64, tags []string) (err error) {
	err = updateMealTags(q, mealID, tags)
	if err != nil {
		return err
	}
	return UpdateSearchText(q, mealID)
}

func updateMealTags(q Queryable, mealID uint64, tags []string) (err error) {
	err = deleteMealTags(q, mealID)
	if err != nil {
		return err
	}

	return addMealTags(q, mealID, tags)
}

// UpdateMealWithTags combines UpdateMeal and UpdateMealTags.
func UpdateMealWithTags(q Queryable, mt mpdata.MealWithTags) (err error) {
	err = updateMealWithTags(q, mt)
	if err != nil {
		return err
	}
	return UpdateSearchText(q, mt.Meal.ID)
}

func updateMealWithTags(q Queryable, mt mpdata.MealWithTags) (err error) {
	err = updateMeal(q, mt.Meal)
	if err != nil {
		return err
	}

	return updateMealTags(q, mt.Meal.ID, mt.Tags)
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

// DeleteMeal deletes the meal record identified by 'mealID'. If no such meal
// exists, no error is raised.
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

	return deleteMealTags(q, mealID)
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
