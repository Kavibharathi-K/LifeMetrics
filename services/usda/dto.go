package usda

type SearchResponse struct {
	Foods []Food `json:"foods"`
}

type Food struct {
	FdcID         int        `json:"fdcId"`
	Description   string     `json:"description"`
	FoodNutrients []Nutrient `json:"foodNutrients"`
}

type Nutrient struct {
	NutrientID int     `json:"nutrientId"`
	Value      float64 `json:"value"`
}