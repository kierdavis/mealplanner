package mpdata

// Type Meal holds information about a meal in the database.
type Meal struct {
	ID        uint64 // The database's unique identifier for the meal.
	Name      string // The name of the meal.
	RecipeURL string // The possibly empty URL of the recipe for the meal.
	Favourite bool   // Whether or not the meal is marked as a favourite.
}

// Type MealWithTags pairs a Meal with its associated tags.
type MealWithTags struct {
	Meal *Meal    // The meal.
	Tags []string // The meal's tags.
}

// Type MealWithScore pairs a Meal with its associated score.
type MealWithScore struct {
	Meal  *Meal   // The meal.
	Score float32 // The meal's score.
}
