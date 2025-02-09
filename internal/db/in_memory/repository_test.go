package in_memory

import (
	"context"
	"github.com/alserok/url_shortener/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestPostgresRepoSuite(t *testing.T) {
	suite.Run(t, new(PostgresRepoSuite))
}

type PostgresRepoSuite struct {
	suite.Suite

	repo *repository
	ctx  context.Context

	mocks struct {
		ctrl *gomock.Controller

		logger *logger.MockLogger
	}
}

func (prs *PostgresRepoSuite) SetupTest() {
	prs.repo = NewRepository()

	prs.mocks.ctrl = gomock.NewController(prs.T())
	prs.mocks.logger = logger.NewMockLogger(prs.mocks.ctrl)

	prs.ctx = logger.WrapLogger(context.Background(), prs.mocks.logger)
}

func (prs *PostgresRepoSuite) TeardownTest() {
	clear(prs.repo.db)
}

func (prs *PostgresRepoSuite) TestSaveURL() {
	url := "url"
	shortened := "u"

	prs.mocks.logger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	prs.Require().NoError(prs.repo.SaveURL(prs.ctx, url, shortened))

	prs.Require().Equal(url, prs.repo.db[shortened].val)
}

func (prs *PostgresRepoSuite) TestGetURL() {
	url := "url"
	shortened := "u"

	prs.mocks.logger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	n := &node[string]{
		key: shortened,
		val: url,
	}
	prs.repo.db[shortened] = n
	prs.repo.head = n
	prs.repo.tail = n

	res, err := prs.repo.GetURL(prs.ctx, shortened)
	prs.Require().NoError(err)
	prs.Require().Equal(url, res)
}

func (prs *PostgresRepoSuite) TestLRULogic() {
	prs.repo.size = 2

	prs.mocks.logger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	// =================

	prs.Require().NoError(prs.repo.SaveURL(prs.ctx, "url", "sh"))
	prs.Require().NoError(prs.repo.SaveURL(prs.ctx, "url1", "sh1"))
	prs.Require().NoError(prs.repo.SaveURL(prs.ctx, "url2", "sh2"))

	prs.Require().Nil(prs.repo.db["sh"])
	prs.Require().Equal(prs.repo.db["sh1"].val, "url1")
	prs.Require().Equal(prs.repo.db["sh2"].val, "url2")
	prs.Require().Equal(prs.repo.head, prs.repo.db["sh2"])
	prs.Require().Equal(prs.repo.tail, prs.repo.db["sh1"])

	// =================

	url, err := prs.repo.GetURL(prs.ctx, "sh1")
	prs.Require().NoError(err)
	prs.Require().Equal(url, "url1")
	prs.Require().Equal(prs.repo.head, prs.repo.db["sh1"])
	prs.Require().Equal(prs.repo.tail, prs.repo.db["sh2"])

	// =================

	prs.Require().NoError(prs.repo.SaveURL(prs.ctx, "url3", "sh3"))
	prs.Require().Nil(prs.repo.db["sh2"])
	prs.Require().Equal(prs.repo.db["sh1"].val, "url1")
	prs.Require().Equal(prs.repo.db["sh3"].val, "url3")
	prs.Require().Equal(prs.repo.head, prs.repo.db["sh3"])
	prs.Require().Equal(prs.repo.tail, prs.repo.db["sh1"])
}

func (prs *PostgresRepoSuite) TestClose() {
	prs.Require().NoError(prs.repo.Close())
}
