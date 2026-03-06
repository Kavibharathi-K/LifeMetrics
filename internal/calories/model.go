package calories

import "github.com/google/uuid"

type Food struct {
	ID uuid.UUID `json:"id"`

	Name            string `json:"name"`
	MeasurementType string `json:"measurement_type"`
	UnitName        string `json:"unit_name"`

	CaloriesPer100g float64 `json:"calories_per_100g"`
	ProteinPer100g  float64 `json:"protein_per_100g"`
	CarbsPer100g    float64 `json:"carbs_per_100g"`
	FatPer100g      float64 `json:"fat_per_100g"`
	FiberPer100g    float64 `json:"fiber_per_100g"`

	CaloriesPerUnit float64 `json:"calories_per_unit"`
	ProteinPerUnit  float64 `json:"protein_per_unit"`
	CarbsPerUnit    float64 `json:"carbs_per_unit"`
	FatPerUnit      float64 `json:"fat_per_unit"`
	FiberPerUnit    float64 `json:"fiber_per_unit"`
}
