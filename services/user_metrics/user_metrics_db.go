package user_metrics

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

func (r *Repository) CreateUserMetrics(
	ctx context.Context,
	m UserMetrics,
) (int, error) {

	query := `
		SELECT user_metrics_add(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10
		)
	`

	var id int

	err := r.DB.QueryRow(
		ctx,
		query,

		m.UserId,

		m.Age,
		m.Gender,

		m.HeightCm,
		m.WeightKg,

		m.ActivityLevel,

		m.MaintenanceCalories,

		m.ProteinGoal,
		m.CarbGoal,
		m.FatGoal,
	).Scan(&id)

	return id, err
}


func (r *Repository) GetLatestUserMetrics(
	ctx context.Context,
	userID int64,
) (UserMetrics, error) {

	var m UserMetrics

	query := `
		SELECT *
		FROM usermetrics_getlatest($1)
	`

	err := r.DB.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&m.Id,
		&m.UserId,

		&m.Age,
		&m.Gender,

		&m.HeightCm,
		&m.WeightKg,

		&m.ActivityLevel,

		&m.MaintenanceCalories,

		&m.ProteinGoal,
		&m.CarbGoal,
		&m.FatGoal,

		&m.CreatedAt,
	)

	return m, err
}