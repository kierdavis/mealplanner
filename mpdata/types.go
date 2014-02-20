package mpdata

import (
	"time"
)

// Meal holds information about a meal in the database.
type Meal struct {
	ID        uint64 `json:"id"`        // The database's unique identifier for the meal.
	Name      string `json:"name"`      // The name of the meal.
	RecipeURL string `json:"recipe"`    // The possibly empty URL of the recipe for the meal.
	Favourite bool   `json:"favourite"` // Whether or not the meal is marked as a favourite.
}

// MealWithTags pairs a Meal with its associated tags.
type MealWithTags struct {
	Meal *Meal    `json:"meal"` // The meal.
	Tags []string `json:"tags"` // The meal's tags.
}

// Suggestion pairs a Meal with its associated tags, closest serving distance and score.
type Suggestion struct {
	MT    MealWithTags `json:"mt"`    // The meal and tags.
	CSD   int          `json:"-"`     // The closest serving distance (used in computing the score).
	Score float32      `json:"score"` // The meal's score.
}

// MealPlan holds information about a meal plan in the database. It
// contains no JSON field tags as the mealPlanJSON struct is actually used for
// encoding/decoding; however, the MarshalJSON/UnmarshalJSON methods take care
// of this.
type MealPlan struct {
	ID        uint64    // The database's unique identifier for the meal plan.
	Notes     string    // The textual notes associated with the meal plan.
	StartDate time.Time // The date of the first day in the meal plan.
	EndDate   time.Time // The date of the last day in the meal plan.
}

// mealPlanJSON is the intermediate struct used for JSON encoding/decoding
// of a meal plan. An intermediate type is used as the time.Times need to be
// encoded in a specific format.
type mealPlanJSON struct {
	ID        uint64 `json:"id"`
	Notes     string `json:"notes"`
	StartDate string `json:"startdate"`
	EndDate   string `json:"enddate"`
}

// Days returns a slice of times representing the days between mp.StartDate and
// mp.EndDate, inclusive.
func (mp *MealPlan) Days() (days []time.Time) {
	curr := mp.StartDate
	for !curr.After(mp.EndDate) {
		days = append(days, curr)
		curr = curr.Add(time.Hour * 24)
	}

	return days
}

// Serving holds information about a serving of a meal in the database.
// It contains no JSON field tags as the servingJSON struct is actually used for
// encoding/decoding; however, the MarshalJSON/UnmarshalJSON methods take care
// of this.
type Serving struct {
	MealPlanID uint64
	Date       time.Time
	MealID     uint64
}

// servingJSON is the intermediate struct used for JSON encoding/decoding
// of a meal plan. An intermediate type is used as the time.Time needs to be
// encoded in a specific format.
type servingJSON struct {
	MealPlanID uint64 `json:"mealplanid"`
	Date       string `json:"date"`
	MealID     uint64 `json:"mealid"`
}

// MealPlanWithServings pairs a MealPlan with its associated Servings.
type MealPlanWithServings struct {
	MealPlan *MealPlan  `json:"mealplan"`
	Servings []*Serving `json:"servings"`
}
