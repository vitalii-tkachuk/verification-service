package cli

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure/config"
)

var (
	ErrDatabaseConnectionFailed = errors.New("database connection failed")
	ErrMigrationFailed          = errors.New("migration failed")
)

const migrationsSourceURL = "file://migrations"
const migrationsDriverName = "postgres"

// RunMigrate open database connection and execute new migrations.
func RunMigrate() error {
	con, err := getConnection()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseConnectionFailed, err)
	}

	defer func() {
		if con != nil {
			_ = con.Close()
		}
	}()

	if err = migrateUp(con); err != nil {
		return fmt.Errorf("%w: %s", ErrMigrationFailed, err)
	}

	fmt.Println("Migration successfully finished.")

	return nil
}

// getConnection parse environment config and open database connection.
func getConnection() (*sql.DB, error) {
	var cfg config.Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	db, err := sql.Open(migrationsDriverName, cfg.PostgresDatabaseDsn())
	if err != nil {
		return nil, err
	}

	return db, nil
}

// migrateUp execute new migrations.
func migrateUp(connection *sql.DB) error {
	driver, err := postgres.WithInstance(connection, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsSourceURL, migrationsDriverName, driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
