package usda

type FoodSuggestion struct {
	Name            string  `json:"name"`
	MeasurementType string  `json:"measurement_type"`
	BaseQuantity    float64 `json:"base_quantity"`

	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Fiber    float64 `json:"fiber"`
}

func MapFood(food Food) FoodSuggestion {
	var suggestion FoodSuggestion

	suggestion.Name = food.Description
	suggestion.MeasurementType = "g"
	suggestion.BaseQuantity = 100

	for _, nutrient := range food.FoodNutrients {
		switch nutrient.NutrientID {
		case 1008, 2047:
			suggestion.Calories = nutrient.Value

		case 1003:
			suggestion.Protein = nutrient.Value

		case 1005:
			suggestion.Carbs = nutrient.Value

		case 1004:
			suggestion.Fat = nutrient.Value

		case 1079:
			suggestion.Fiber = nutrient.Value
		}
	}

	return suggestion
}