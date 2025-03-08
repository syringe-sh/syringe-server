package database

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed sql/*_tenant.*.sql
var TenantMigrations embed.FS

//go:embed sql/*_system.*.sql
var SystemMigrations embed.FS

type Migrator interface {
	Up() error
	Down() error
}

type Migration struct {
	migrate *migrate.Migrate
}

func (m *Migration) Up() error {
	return m.migrate.Up()
}

func (m *Migration) Down() error {
	return m.migrate.Down()
}

func NewMigration(db *sql.DB, migrations source.Driver) (*Migration, error) {
	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return nil, fmt.Errorf("sqlite3 instance: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"file",
		migrations,
		"sqlite3",
		instance,
	)
	if err != nil {
		return nil, fmt.Errorf("create migration: %w", err)
	}

	return &Migration{migrate: m}, nil
}
