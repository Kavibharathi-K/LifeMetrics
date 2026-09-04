package meal_item

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) AddMealItem(
	ctx context.Context,
	userID int64,
	mealType string,
	mealDate string,
	foodName string,
	quantity float64,
) error {

	query := `
		CALL add_meal_item(
			$1,
			$2,
			$3,
			$4,
			$5
		)
	`

	_, err := r.DB.Exec(
		ctx,
		query,
		userID,
		mealType,
		mealDate,
		foodName,
		quantity,
	)

	return err
}

func (r *Repository) GetMealItems(
	ctx context.Context,
	userID int64,
) ([]MealItem, error) {

	query := `
		SELECT
			mi.meal_item_id,
			mi.meal_id,
			m.meal_type,
			mi.food_id,
			f.name AS food_name,
			mi.quantity,
			mi.calories,
			mi.protein,
			mi.carbs,
			mi.fat,
			mi.fiber,
			mi.created_at
		FROM meal_items mi
		JOIN meals m
			ON m.meal_id = mi.meal_id
		JOIN foods f
			ON f.food_id = mi.food_id
		WHERE m.user_id = $1
		ORDER BY mi.created_at DESC
	`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []MealItem

	for rows.Next() {

		var item MealItem

		err := rows.Scan(
			&item.MealItemID,
			&item.MealID,
			&item.MealType,
			&item.FoodID,
			&item.FoodName,
			&item.Quantity,
			&item.Calories,
			&item.Protein,
			&item.Carbs,
			&item.Fat,
			&item.Fiber,
			&item.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}