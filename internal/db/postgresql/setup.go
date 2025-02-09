package postgresql

import (
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"

	_ "github.com/lib/pq"
)

const (
	postgresDialect = "postgres"
)

func MustConnect(dsn, migrationsDir string) *sqlx.DB {
	db, err := sqlx.Connect(postgresDialect, dsn)
	if err != nil {
		panic("failed to connect to postgres: " + err.Error())
	}

	if err = db.Ping(); err != nil {
		panic("failed to ping postgres: " + err.Error())
	}

	mustMigrate(db, migrationsDir)

	return db
}

func mustMigrate(db *sqlx.DB, migrationsDir string) {
	if err := goose.SetDialect(postgresDialect); err != nil {
		panic(err)
	}

	if err := goose.Up(db.DB, migrationsDir); err != nil {
		panic(err)
	}
}
