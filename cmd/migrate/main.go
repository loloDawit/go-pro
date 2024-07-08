package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"github.com/urfave/cli/v2"

	"github.com/loloDawit/ecom/config"
)

func main() {
	app := &cli.App{
		Name:  "db-migrate",
		Usage: "Run database migrations",
		Commands: []*cli.Command{
			{
				Name:   "up",
				Usage:  "Apply all up migrations",
				Action: runMigrationsUp,
			},
			{
				Name:   "down",
				Usage:  "Revert the last migration",
				Action: runMigrationsDown,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runMigrationsUp(c *cli.Context) error {
	return runMigrations("up")
}

func runMigrationsDown(c *cli.Context) error {
	return runMigrations("down")
}

func runMigrations(direction string) error {
	ctx := context.Background()

	// Set the directory, environment, and deployment
	directory := os.Getenv("CONFIG_DIRECTORY")
	if directory == "" {
		directory = "./config"
	}
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "development"
	}
	deployment := os.Getenv("DEPLOYMENT")

	cfg := config.LoadConfig(ctx, directory, environment, deployment)
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=require", cfg.DBuser, cfg.DBpassword, cfg.DBaddr, cfg.DBname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("unable to connect to the database: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create the PostgreSQL driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create the migrate instance: %v", err)
	}

	switch direction {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	default:
		return fmt.Errorf("unknown migration direction: %s", direction)
	}

	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %v", err)
	}

	fmt.Printf("Migrations %s applied successfully\n", direction)
	return nil
}
