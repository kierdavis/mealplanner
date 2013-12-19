package mpdata

type Meal struct {
	ID uint64
	Name string
	RecipeURL string
	Favourite bool
}

type MealScore struct {
	Meal *Meal
	Score float32
}
