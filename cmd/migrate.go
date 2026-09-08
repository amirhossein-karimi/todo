package cmd

import (
	"fmt"

	"github.com/amirhossein-karimi/todo/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

func createMigrator(cfg *config.Config) (*migrate.Migrate, error) {
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}

	return m, nil
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("config.yaml")
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		m, err := createMigrator(cfg)
		if err != nil {
			return err
		}

		defer m.Close()

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migration up: %w", err)
		}

		fmt.Println("Migrations applied successfully")

		return nil
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("config.yaml")
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		m, err := createMigrator(cfg)
		if err != nil {
			return err
		}

		defer m.Close()

		if err := m.Steps(-1); err != nil {
			return fmt.Errorf("migration down: %w", err)
		}

		fmt.Println("Migration rolled back successfully")

		return nil
	},
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)

	rootCmd.AddCommand(migrateCmd)
}
