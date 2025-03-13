package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*
var migrationsFs embed.FS

var migrationOnce sync.Once
var migrationErr error

// MigrateOnce выполняет миграции базы данных.
func MigrateOnce(db *sql.DB) error {
	migrationOnce.Do(func() {
		migrationErr = Migrate(db)
	})

	return migrationErr
}

func Migrate(db *sql.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create Sqlite migration driver: %w", err)
	}

	sourceDriver, err := iofs.New(migrationsFs, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create source driver for embedded migrations: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		sourceDriver,
		"wb-go",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("no new migrations to apply")
	} else {
		log.Println("migrations applied successfully")
	}

	return nil
}
