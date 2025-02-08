package postgresql

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func TestPostgresRepoSuite(t *testing.T) {
	suite.Run(t, new(PostgresRepoSuite))
}

type PostgresRepoSuite struct {
	suite.Suite

	conn *sqlx.DB
	repo *repository

	containers struct {
		postgres *postgres.PostgresContainer
	}
}

func (prs *PostgresRepoSuite) SetupTest() {
	prs.newPostgresDB()
	prs.repo = NewRepository(prs.conn)
}

func (prs *PostgresRepoSuite) TeardownTest() {
	prs.Require().NoError(prs.containers.postgres.Terminate(context.Background()))
}

func (prs *PostgresRepoSuite) TestSaveURL() {
	url := "url"
	shortened := "u"

	prs.Require().NoError(prs.repo.SaveURL(context.Background(), url, shortened))

	var res string
	prs.Require().NoError(prs.conn.QueryRowx(`SELECT url FROM urls WHERE shortened_url = $1`, shortened).Scan(&res))
	prs.Require().Equal(url, res)
}

func (prs *PostgresRepoSuite) TestGetURL() {
	url := "url"
	shortened := "u"

	prs.Require().NoError(prs.conn.QueryRowx(`INSERT INTO urls (url, shortened_url) VALUES ($1,$2)`, url, shortened).Err())

	res, err := prs.repo.GetURL(context.Background(), shortened)
	prs.Require().NoError(err)
	prs.Require().Equal(url, res)
}

func (prs *PostgresRepoSuite) TestClose() {
	prs.Require().NoError(prs.repo.Close())
}

func (prs *PostgresRepoSuite) newPostgresDB() {
	ctx := context.Background()

	var (
		name     = "postgres"
		user     = "postgres"
		password = "postgres"
	)

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(name),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	prs.Require().NoError(err)
	prs.Require().NotNil(postgresContainer)
	prs.Require().True(postgresContainer.IsRunning())

	port, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	prs.Require().NoError(err)

	conn := MustConnect(
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, "localhost", port.Port(), name),
		"./migrations",
	)

	prs.containers.postgres, prs.conn = postgresContainer, conn
}
