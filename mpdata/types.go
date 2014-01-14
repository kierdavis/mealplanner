package mpdata

import (
	"time"
)

// Type Meal holds information about a meal in the database.
type Meal struct {
	ID        uint64 `json:"id"`        // The database's unique identifier for the meal.
	Name      string `json:"name"`      // The name of the meal.
	RecipeURL string `json:"recipe"`    // The possibly empty URL of the recipe for the meal.
	Favourite bool   `json:"favourite"` // Whether or not the meal is marked as a favourite.
}

// Type MealWithTags pairs a Meal with its associated tags.
type MealWithTags struct {
	Meal *Meal    `json:"meal"` // The meal.
	Tags []string `json:"tags"` // The meal's tags.
}

// Type MealWithScore pairs a Meal with its associated score.
type MealWithScore struct {
	Meal  *Meal   `json:"meal"`  // The meal.
	Score float32 `json:"score"` // The meal's score.
}

type MealPlan struct {
	ID        uint64    `json:"id"`
	Notes     string    `json:"notes"`
	StartDate time.Time `json:"start"`
	EndDate   time.Time `json:"end"`
}

func (mp *MealPlan) Days() (days []time.Time) {
	curr := mp.StartDate
	for !curr.After(mp.EndDate) {
		days = append(days, curr)
		curr = curr.Add(time.Hour * 24)
	}

	return days
}

type Serving struct {
	MealPlanID uint64    `json:"mealplanid"`
	Date       time.Time `json:"date"`
	MealID     uint64    `json:"mealid"`
}

type MealPlanWithServings struct {
	MealPlan *MealPlan  `json:"mealplan"`
	Servings []*Serving `json:"servings"`
}
