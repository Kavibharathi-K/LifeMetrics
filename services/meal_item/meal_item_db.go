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
	mealType string,
	mealDate string,
	foodName string,
	quantity float64,
) error {

	query := `
		call add_meal_item(
			$1,
			$2,
			$3,
			$4
		)
	`
	_, err := r.DB.Exec(
		ctx,
		query,
		mealType,
		mealDate,
		foodName,
		quantity,
	)
	return err
}
