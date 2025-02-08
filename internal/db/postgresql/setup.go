package postgresql

import (
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"

	_ "github.com/lib/pq"
)

func MustConnect(dsn, migrationsDir string) *sqlx.DB {
	db, err := sqlx.Connect(postgresDialect, dsn)
	if err != nil {
		panic("failed to connect to postgres: " + err.Error())
	}

	mustMigrate(db, migrationsDir)

	return db
}

const (
	postgresDialect = "postgres"
)

func mustMigrate(db *sqlx.DB, migrationsDir string) {
	if err := goose.SetDialect(postgresDialect); err != nil {
		panic(err)
	}

	if err := goose.Up(db.DB, migrationsDir); err != nil {
		panic(err)
	}
}
