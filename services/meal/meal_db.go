package meal

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetTodayNutrition(ctx context.Context) ([]TodayNutrition, error) {
	today := time.Now().Format("2006-01-02")

	query := `select * from get_today_nutrition($1)`
	rows, err := r.DB.Query(ctx, query, today)
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
