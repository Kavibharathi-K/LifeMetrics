package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(pool *pgxpool.Pool) error {

	ctx := context.Background()

	_, err := pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT NOW()
	);
	`)
	if err != nil {
		return err
	}

	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		return err
	}

	log.Println("Files", files)

	sort.Strings(files)

	for _, file := range files {

		version := filepath.Base(file)

		var exists bool

		err := pool.QueryRow(
			ctx,
			"SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version=$1)",
			version,
		).Scan(&exists)

		if err != nil {
			return err
		}

		if exists {
			continue
		}

		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		_, err = pool.Exec(ctx, string(sqlBytes))
		if err != nil {
			return err
		}

		_, err = pool.Exec(
			ctx,
			"INSERT INTO schema_migrations (version) VALUES ($1)",
			version,
		)

		if err != nil {
			return err
		}

		fmt.Println("✅ Applied migration:", version)
	}

	return nil
}
