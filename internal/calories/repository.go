package calories

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

	query := `
	INSERT INTO foods (
		name,
		measurement_type,
		unit_name,
		calories_per_100g,
		protein_per_100g,
		carbs_per_100g,
		fat_per_100g,
		fiber_per_100g,
		calories_per_unit,
		protein_per_unit,
		carbs_per_unit,
		fat_per_unit,
		fiber_per_unit
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	RETURNING id;
	`

	return r.DB.QueryRow(ctx, query,
		food.Name,
		food.MeasurementType,
		food.UnitName,
		food.CaloriesPer100g,
		food.ProteinPer100g,
		food.CarbsPer100g,
		food.FatPer100g,
		food.FiberPer100g,
		food.CaloriesPerUnit,
		food.ProteinPerUnit,
		food.CarbsPerUnit,
		food.FatPerUnit,
		food.FiberPerUnit,
	).Scan(&food.ID)
}

func (r *Repository) ListFoods(ctx context.Context) ([]Food, error) {

	rows, err := r.DB.Query(ctx, `
	SELECT 
	id,
	name,
	measurement_type,
	unit_name,
	calories_per_100g,
	protein_per_100g,
	carbs_per_100g,
	fat_per_100g,
	fiber_per_100g,
	calories_per_unit,
	protein_per_unit,
	carbs_per_unit,
	fat_per_unit,
	fiber_per_unit
	FROM foods`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	foods := make([]Food, 0)

	for rows.Next() {

		var f Food

		err := rows.Scan(
			&f.ID,
			&f.Name,
			&f.MeasurementType,
			&f.UnitName,
			&f.CaloriesPer100g,
			&f.ProteinPer100g,
			&f.CarbsPer100g,
			&f.FatPer100g,
			&f.FiberPer100g,
			&f.CaloriesPerUnit,
			&f.ProteinPerUnit,
			&f.CarbsPerUnit,
			&f.FatPerUnit,
			&f.FiberPerUnit,
		)

		if err != nil {
			return nil, err
		}

		foods = append(foods, f)
	}

	return foods, nil
}