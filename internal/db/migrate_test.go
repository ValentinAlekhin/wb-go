package db

import (
	"database/sql"
	"embed"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// resetGlobals сбрасывает глобальные переменные для изоляции тестов.
func resetGlobals() {
	migrationOnce = sync.Once{}
	migrationErr = nil
}

// TestRunMigrations тестирует функцию MigrateOnce.
func TestRunMigrations(t *testing.T) {
	tests := []struct {
		name        string
		db          *sql.DB
		migrationFS embed.FS
		wantErr     bool
		errContains string
		validate    func(t *testing.T, db *sql.DB)
	}{
		{
			name:        "nil database",
			db:          nil,
			migrationFS: migrationsFs,
			wantErr:     true,
			errContains: "db is nil",
		},
		{
			name: "successful migration",
			db: func() *sql.DB {
				db, err := sql.Open("sqlite", ":memory:")
				require.NoError(t, err, "Failed to create in-memory DB")
				return db
			}(),
			migrationFS: migrationsFs,
			wantErr:     false,
			validate: func(t *testing.T, db *sql.DB) {
				// Проверяем, что таблица virtual_controls создана
				rows, err := db.Query("SELECT topic, value, created_at, updated_at FROM virtual_controls")
				require.NoError(t, err, "Failed to query virtual_controls")
				defer rows.Close()
				assert.False(t, rows.Next(), "Table should be empty initially")
			},
		},
		{
			name: "no migrations to apply",
			db: func() *sql.DB {
				db, err := sql.Open("sqlite", ":memory:")
				require.NoError(t, err, "Failed to create in-memory DB")
				// Применяем миграции заранее
				applyMigrations(t, db)
				return db
			}(),
			migrationFS: migrationsFs,
			wantErr:     false,
			validate: func(t *testing.T, db *sql.DB) {
				// Проверяем, что таблица всё ещё существует
				_, err := db.Exec("INSERT INTO virtual_controls (topic, value, created_at, updated_at) VALUES (?, ?, ?, ?)",
					"test/topic", "123", time.Now(), time.Now())
				require.NoError(t, err, "Failed to insert into virtual_controls")
			},
		},
		{
			name: "empty migrations directory",
			db: func() *sql.DB {
				db, err := sql.Open("sqlite", ":memory:")
				require.NoError(t, err, "Failed to create in-memory DB")
				return db
			}(),
			migrationFS: embed.FS{}, // Пустая FS
			wantErr:     true,
			errContains: "failed to create source driver",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сбрасываем глобальные переменные
			resetGlobals()
			migrationsFs = tt.migrationFS

			// Выполняем миграции
			err := MigrateOnce(tt.db)

			// Проверяем результат с использованием testify
			if tt.wantErr {
				assert.Error(t, err, "MigrateOnce() should return an error")
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains, "Error message should contain expected substring")
				}
			} else {
				assert.NoError(t, err, "MigrateOnce() should not return an error")
				if tt.validate != nil {
					tt.validate(t, tt.db)
				}
			}

			// Закрываем DB, если она была открыта
			if tt.db != nil {
				tt.db.Close()
			}
		})
	}
}

// applyMigrations применяет миграции вручную для подготовки теста.
func applyMigrations(t *testing.T, db *sql.DB) {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	require.NoError(t, err, "Failed to create SQLite driver")

	sourceDriver, err := iofs.New(migrationsFs, "migrations")
	require.NoError(t, err, "Failed to create source driver")

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "wb-go", driver)
	require.NoError(t, err, "Failed to initialize migrations")

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		t.Fatalf("Failed to apply migrations: %v", err)
	}
}
