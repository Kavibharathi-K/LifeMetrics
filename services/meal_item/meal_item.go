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
