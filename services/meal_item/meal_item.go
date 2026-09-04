package meal_item

import "time"

type MealItem struct {
	MealItemID int `db:"meal_item_id" json:"meal_item_id"`

	MealID int `db:"meal_id" json:"meal_id"`
	MealType string `db:"meal_type" json:"meal_type"`
	FoodID int `db:"food_id" json:"food_id"`
	FoodName string `db:"food_name" json:"food_name"`

	Quantity float64 `db:"quantity" json:"quantity"`

	Calories float64 `db:"calories" json:"calories"`
	Protein  float64 `db:"protein" json:"protein"`
	Carbs    float64 `db:"carbs" json:"carbs"`
	Fat      float64 `db:"fat" json:"fat"`
	Fiber    float64 `db:"fiber" json:"fiber"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
}


type MealItemResponse struct {
	MealID        int     `json:"meal_id"`
	MealType      string  `json:"meal_type"`
	TotalCalories float64 `json:"total_calories"`
	TotalProtein  float64 `json:"total_protein"`
	TotalCarbs    float64 `json:"total_carbs"`
	TotalFat      float64 `json:"total_fat"`
	TotalFiber    float64 `json:"total_fiber"`

	FoodID          int     `json:"food_id"`
	FoodName        string  `json:"food_name"`
	MeasurementType string  `json:"measurement_type"`

	Quantity float64 `json:"quantity"`
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Fiber    float64 `json:"fiber"`
}