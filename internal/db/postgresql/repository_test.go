package postgresql

import (
	"context"
	"fmt"
	"github.com/alserok/url_shortener/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

func TestPostgresRepoSuite(t *testing.T) {
	suite.Run(t, new(PostgresRepoSuite))
}

type PostgresRepoSuite struct {
	suite.Suite

	conn *sqlx.DB
	repo *repository
	ctx  context.Context

	containers struct {
		postgres *postgres.PostgresContainer
	}

	mocks struct {
		ctrl *gomock.Controller

		logger *logger.MockLogger
	}
}

func (prs *PostgresRepoSuite) SetupTest() {
	prs.newPostgresDB()
	prs.repo = NewRepository(prs.conn)

	prs.mocks.ctrl = gomock.NewController(prs.T())
	prs.mocks.logger = logger.NewMockLogger(prs.mocks.ctrl)

	prs.ctx = logger.WrapLogger(context.Background(), prs.mocks.logger)
}

func (prs *PostgresRepoSuite) TeardownTest() {
	prs.Require().NoError(prs.containers.postgres.Terminate(context.Background()))
}

func (prs *PostgresRepoSuite) TestSaveURL() {
	url := "url"
	shortened := "u"

	prs.mocks.logger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	prs.Require().NoError(prs.repo.SaveURL(prs.ctx, url, shortened))

	var res string
	prs.Require().NoError(prs.conn.QueryRowx(`SELECT url FROM urls WHERE shortened_url = $1`, shortened).Scan(&res))
	prs.Require().Equal(url, res)
}

func (prs *PostgresRepoSuite) TestGetURL() {
	url := "url"
	shortened := "u"

	prs.mocks.logger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	prs.Require().NoError(prs.conn.QueryRowx(`INSERT INTO urls (url, shortened_url) VALUES ($1,$2)`, url, shortened).Err())

	res, err := prs.repo.GetURL(prs.ctx, shortened)
	prs.Require().NoError(err)
	prs.Require().Equal(url, res)
}

func (prs *PostgresRepoSuite) TestClose() {
	prs.Require().NoError(prs.repo.Close())
}

// launches postgres container and connect to it
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
				WithOccurrence(2),
		),
	)
	prs.Require().NoError(err)
	prs.Require().NotNil(postgresContainer)
	prs.Require().True(postgresContainer.IsRunning())

	port, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	prs.Require().NoError(err)

	conn := MustConnect(
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, "localhost", port.Port(), name),
		"../../../migrations/postgres",
	)

	prs.containers.postgres, prs.conn = postgresContainer, conn
}
