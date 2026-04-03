package food

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

func (r *Repository) CreateFood(ctx context.Context, food *Food) error {

	query := `SELECT create_food($1,$2,$3,$4,$5,$6,$7,$8);`

	return r.DB.QueryRow(ctx, query,
		food.Name,
		food.MeasurementType,
		food.BaseQuantity,
		food.Calories,
		food.Protein,
		food.Carbs,
		food.Fat,
		food.Fiber,
	).Scan(&food.FoodID)
}

func (r *Repository) ListFoods(ctx context.Context) ([]Food, error) {

	rows, err := r.DB.Query(ctx, `SELECT * FROM list_foods();`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	foods := make([]Food, 0)

	for rows.Next() {

		var f Food

		err := rows.Scan(
			&f.FoodID,
			&f.Name,
			&f.MeasurementType,
			&f.BaseQuantity,
			&f.Calories,
			&f.Protein,
			&f.Carbs,
			&f.Fat,
			&f.Fiber,
			&f.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		foods = append(foods, f)
	}

	return foods, nil
}

func (r *Repository) GetFoodByID(ctx context.Context, foodID int,) (Food, error) {

	query := `
	SELECT * FROM get_food_by_id($1)
	`
	
	var food Food

	err := r.DB.QueryRow(
		ctx,
		query,
		foodID,
	).Scan(
		&food.FoodID,
		&food.Name,
		&food.MeasurementType,
		&food.BaseQuantity,
		&food.Calories,
		&food.Protein,
		&food.Carbs,
		&food.Fat,
		&food.Fiber,
		&food.CreatedAt,
	)

	return food, err
}

func (r *Repository) UpdateFood(ctx context.Context, foodID int, food *Food) error {
	query := `
	call update_food($1,$2,$3,$4,$5,$6,$7,$8,$9);
	`

	_, err := r.DB.Exec(
		ctx, 
		query,
		foodID,
		food.Name,
		food.MeasurementType,
		food.BaseQuantity,
		food.Calories,
		food.Protein,
		food.Carbs,
		food.Fat,
		food.Fiber,
	)
	return err
}
