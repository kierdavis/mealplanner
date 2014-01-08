package mpdata

// Type Meal holds information about a meal in the database.
type Meal struct {
	ID        uint64 `json:"id"` // The database's unique identifier for the meal.
	Name      string `json:"name"` // The name of the meal.
	RecipeURL string `json:"recipe"` // The possibly empty URL of the recipe for the meal.
	Favourite bool   `json:"favourite"` // Whether or not the meal is marked as a favourite.
}

// Type MealWithTags pairs a Meal with its associated tags.
type MealWithTags struct {
	Meal *Meal    `json:"meal"` // The meal.
	Tags []string `json:"tags"` // The meal's tags.
}

// Type MealWithScore pairs a Meal with its associated score.
type MealWithScore struct {
	Meal  *Meal   `json:"meal"` // The meal.
	Score float32 `json:"score"` // The meal's score.
}

type MealPlan struct {
	ID uint64 `json:"id"`
	Notes string `json:"notes"`
	StartDate time.Time `json:"start"`
	EndDate time.Time `json:"end"`
}
