package food

import "time"

type Food struct {
	FoodID int    `db:"food_id" json:"food_id"`
	Name   string `db:"name" json:"name"`

	MeasurementType string  `db:"measurement_type" json:"measurement_type"`
	BaseQuantity    float64 `db:"base_quantity" json:"base_quantity"`

	Calories float64 `db:"calories" json:"calories"`
	Protein  float64 `db:"protein" json:"protein"`
	Carbs    float64 `db:"carbs" json:"carbs"`
	Fat      float64 `db:"fat" json:"fat"`
	Fiber    float64 `db:"fiber" json:"fiber"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
