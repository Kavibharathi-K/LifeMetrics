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

	log.Println("Running TABLE migrations...")
	//Table migrations
	err := runFolderMigrations(ctx, pool, "migrations/tables")
	if err != nil {
		return err
	}

	log.Println("Running PROC migrations...")
	//Proc migrations
	err = runFolderMigrations(ctx, pool, "migrations/procs")
	if err != nil {
		return err
	}

	log.Println("All migrations completed ✅")

	return nil
}

func runFolderMigrations(ctx context.Context, pool *pgxpool.Pool, folder string) error {

	files, err := filepath.Glob(filepath.Join(folder, "*.sql"))
	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, file := range files {

		version := filepath.Base(file)

		var exists bool

		err := pool.QueryRow(
			ctx,
			`SELECT EXISTS (
				SELECT 1
				FROM schema_migrations
				WHERE version=$1
			)`,
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

		log.Println("Applying:", version)

		_, err = pool.Exec(ctx, string(sqlBytes))
		if err != nil {
			return err
		}

		_, err = pool.Exec(
			ctx,
			`INSERT INTO schema_migrations (version)
			 VALUES ($1)`,
			version,
		)
		if err != nil {
			return err
		}

		fmt.Println("✅ Applied migration:", version)
	}

	return nil
}
