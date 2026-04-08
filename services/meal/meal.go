package meal

import "time"

type Meal struct {
	MealID int `db:"meal_id" json:"meal_id"`

	MealDate time.Time `db:"meal_date" json:"meal_date"`

	MealType string `db:"meal_type" json:"meal_type"`
	// breakfast
	// lunch
	// dinner
	// snack

	TotalCalories float64 `db:"total_calories" json:"total_calories"`
	TotalProtein  float64 `db:"total_protein" json:"total_protein"`
	TotalCarbs    float64 `db:"total_carbs" json:"total_carbs"`
	TotalFat      float64 `db:"total_fat" json:"total_fat"`
	TotalFiber    float64 `db:"total_fiber" json:"total_fiber"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type TodayNutrition struct {
	TotalCalories float64 `db:"total_calories" json:"total_calories"`
	TotalProtein  float64 `db:"total_protein" json:"total_protein"`
	TotalCarbs    float64 `db:"total_carbs" json:"total_carbs"`
	TotalFat      float64 `db:"total_fat" json:"total_fat"`
	TotalFiber    float64 `db:"total_fiber" json:"total_fiber"`
}
