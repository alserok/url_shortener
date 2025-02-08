package in_memory

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestPostgresRepoSuite(t *testing.T) {
	suite.Run(t, new(PostgresRepoSuite))
}

type PostgresRepoSuite struct {
	suite.Suite

	repo *repository
}

func (prs *PostgresRepoSuite) SetupTest() {
	prs.repo = NewRepository()
}

func (prs *PostgresRepoSuite) TeardownTest() {
	clear(prs.repo.db)
}

func (prs *PostgresRepoSuite) TestSaveURL() {
	url := "url"
	shortened := "u"

	prs.Require().NoError(prs.repo.SaveURL(context.Background(), url, shortened))

	prs.Require().Equal(url, prs.repo.db[shortened])
}

func (prs *PostgresRepoSuite) TestGetURL() {
	url := "url"
	shortened := "u"

	prs.repo.db[shortened] = url

	res, err := prs.repo.GetURL(context.Background(), shortened)
	prs.Require().NoError(err)
	prs.Require().Equal(url, res)
}

func (prs *PostgresRepoSuite) TestClose() {
	prs.Require().NoError(prs.repo.Close())
}
