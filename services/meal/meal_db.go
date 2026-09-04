package meal

import (
	"context"
	"math"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetTodayNutrition(ctx context.Context, userID int64) ([]TodayNutrition, error) {
	today := time.Now().Format("2006-01-02")

	query := `select * from get_today_nutrition($1, $2)`
	rows, err := r.DB.Query(ctx, query, userID, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []TodayNutrition
	for rows.Next() {
		var m TodayNutrition
		err := rows.Scan(
			&m.TotalCalories,
			&m.TotalProtein,
			&m.TotalCarbs,
			&m.TotalFat,
			&m.TotalFiber,
		)
		if err != nil {
			return nil, err
		}
		meals = append(meals, m)
	}
	return meals, nil

}

func (r *Repository) GetMealsByDate(ctx context.Context, userID int64, date string) (MealsResponse, error) {

	rows, err := r.DB.Query(ctx,
		`SELECT * FROM get_meals_by_date($1, $2)`,
		userID,
		date,
	)

	if err != nil {
		return MealsResponse{}, err
	}

	defer rows.Close()

	response := MealsResponse{
		Breakfast: Meals{Foods: []FoodItem{}},
		Lunch:     Meals{Foods: []FoodItem{}},
		Dinner:    Meals{Foods: []FoodItem{}},
		Snack:     Meals{Foods: []FoodItem{}},
	}

	for rows.Next() {

		var mealType string
		var mealID int
		var totals MealTotals
		var item FoodItem

		err := rows.Scan(

			&mealID,
			&mealType,

			&totals.Calories,
			&totals.Protein,
			&totals.Carbs,
			&totals.Fat,
			&totals.Fiber,

			&item.FoodID,
			&item.Name,
			&item.Unit,

			&item.Quantity,
			&item.Calories,
			&item.Protein,
			&item.Carbs,
			&item.Fat,
			&item.Fiber,
		)

		if err != nil {
			return MealsResponse{}, err
		}

		var meal *Meals

		switch mealType {

		case "breakfast":
			meal = &response.Breakfast

		case "lunch":
			meal = &response.Lunch

		case "dinner":
			meal = &response.Dinner

		case "snack", "snacks":
			meal = &response.Snack

		default:
			continue
		}

		if meal == nil {
			continue
		}

		if meal.MealID == 0 {

			meal.MealID = mealID
			meal.Totals = totals
		}

		if item.FoodID != 0 {

			meal.Foods =
				append(meal.Foods, item)
		}
	}

	addTotals(
		&response.DailyTotals,
		response.Breakfast.Totals,
	)

	addTotals(
		&response.DailyTotals,
		response.Lunch.Totals,
	)

	addTotals(
		&response.DailyTotals,
		response.Dinner.Totals,
	)

	addTotals(
		&response.DailyTotals,
		response.Snack.Totals,
	)

	return response, nil
}

func round(val float64) float64 {
	return math.Round(val*100) / 100
}

func addTotals(total *MealTotals, meal MealTotals) {

	total.Calories = round(total.Calories + meal.Calories)
	total.Protein = round(total.Protein + meal.Protein)
	total.Carbs = round(total.Carbs + meal.Carbs)
	total.Fat = round(total.Fat + meal.Fat)
	total.Fiber = round(total.Fiber + meal.Fiber)
}
