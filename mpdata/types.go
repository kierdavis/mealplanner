package mpdata

type Meal struct {
	ID uint64
	Name string
	RecipeURL string
	Favourite bool
	Tags []string
	HasTags bool
}

type MealWithScore struct {
	Meal *Meal
	Score float32
}
