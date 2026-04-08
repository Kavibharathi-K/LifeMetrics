package user_metrics

import "time"

var activityMultiplier = map[string]float64{
	"sedentary":         1.2,
	"lightly_active":    1.375,
	"moderately_active": 1.55,
	"very_active":       1.725,
	"extra_active":      1.9,
}

type UserMetrics struct {
	Id int `db:"id" json:"id"`

	Age    int    `db:"age" json:"age"`
	Gender string `db:"gender" json:"gender"`

	HeightCm float64 `db:"height_cm" json:"height_cm"`
	WeightKg float64 `db:"weight_kg" json:"weight_kg"`

	ActivityLevel string `db:"activity_level" json:"activity_level"`

	MaintenanceCalories int `db:"maintenance_calories" json:"maintenance_calories"`

	ProteinGoal int `json:"protein_goal"`
	CarbGoal    int `json:"carb_goal"`
	FatGoal     int `json:"fat_goal"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func CalculateMaintenanceCalories(age int, gender string, heightCm float64, weightKg float64, activityLevel string) int {

	var bmr float64

	if gender == "male" {
		bmr =
			10*weightKg +
				6.25*heightCm -
				5*float64(age) +
				5

	} else {

		bmr =
			10*weightKg +
				6.25*heightCm -
				5*float64(age) -
				161
	}

	multiplier := map[string]float64{

		"sedentary":         1.2,
		"lightly_active":    1.375,
		"moderately_active": 1.55,
		"very_active":       1.725,
		"extra_active":      1.9,
	}

	return int(bmr * multiplier[activityLevel])
}

func CalculateMacros(calories int, proteinRatio float64, carbRatio float64, fatRatio float64) (int, int, int) {

	proteinGrams :=
		int((float64(calories) * proteinRatio / 100) / 4)

	carbGrams :=
		int((float64(calories) * carbRatio / 100) / 4)

	fatGrams :=
		int((float64(calories) * fatRatio / 100) / 9)

	return proteinGrams, carbGrams, fatGrams
}

func CalculateMacroGoals(
	calories int,
	weightKg int,
) (int, int, int) {

	// protein based on body weight
	proteinGrams :=
		int(float64(weightKg) * 2.0)

	// fats fixed to 25% calories
	fatCalories :=
		float64(calories) * 0.25

	fatGrams :=
		int(fatCalories / 9)

	// remaining calories go to carbs
	carbCalories :=
		float64(calories) * 0.45

	carbGrams :=
		int(carbCalories / 4)

	return proteinGrams,
		carbGrams,
		fatGrams
}
