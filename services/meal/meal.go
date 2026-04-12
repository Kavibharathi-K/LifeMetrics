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

type FoodItem struct {
	FoodID   int     `json:"food_id"`
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`

	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Fiber    float64 `json:"fiber"`
}

type MealTotals struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Fiber    float64 `json:"fiber"`
}

type Meals struct {
	MealID int        `json:"meal_id"`
	Foods  []FoodItem `json:"foods"`
	Totals MealTotals `json:"totals"`
}

type MealsResponse struct {
	Breakfast   Meals      `json:"breakfast"`
	Lunch       Meals      `json:"lunch"`
	Dinner      Meals       `json:"dinner"`
	Snack       Meals       `json:"snack"`
	DailyTotals MealTotals `json:"daily_totals"`
}
