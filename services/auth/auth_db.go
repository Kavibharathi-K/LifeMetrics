package auth

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

func (r *Repository) CreateUser(ctx context.Context, email string, passwordHash string) (int64, error) {
	query := `
		SELECT user_add(
			$1,
			$2
		)
	`
	var id int64
	err := r.DB.QueryRow(
		ctx,
		query,
		email,
		passwordHash,
	).Scan(&id)
	return id, err
}

func (r *Repository) GetUserByEmail(ctx context.Context,email string) (User, error) {
	query := `
		SELECT *
		FROM user_get_by_email($1)
	`
	var user User
	err := r.DB.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}