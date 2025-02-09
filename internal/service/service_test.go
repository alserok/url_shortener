package service

import (
	"context"
	"github.com/alserok/url_shortener/internal/db"
	"github.com/alserok/url_shortener/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	mocks struct {
		ctrl *gomock.Controller

		repo   *db.MockRepository
		logger *logger.MockLogger
	}
}

func (ss *ServiceSuite) SetupTest() {
	ss.mocks.ctrl = gomock.NewController(ss.T())
	ss.mocks.repo = db.NewMockRepository(ss.mocks.ctrl)
	ss.mocks.logger = logger.NewMockLogger(ss.mocks.ctrl)

	ss.ctx = logger.WrapLogger(context.Background(), ss.mocks.logger)
}

func (ss *ServiceSuite) TeardownTest() {
	ss.mocks.ctrl.Finish()
}

func (ss *ServiceSuite) TestShortenAndSaveURL() {
	url := "https://myurl.com"

	ss.mocks.repo.EXPECT().
		SaveURL(gomock.Any(), gomock.Eq(url), gomock.Any()).
		Return(nil).
		Times(1)

	ss.mocks.logger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	res, err := New(ss.mocks.repo).ShortenAndSaveURL(ss.ctx, url)
	ss.Require().NoError(err)
	ss.Require().NotEqual("", res)
}

func (ss *ServiceSuite) TestGetURL() {
	shortened := "http://c46ce31003"
	url := "http://my_addr.com/"

	ss.mocks.repo.EXPECT().
		GetURL(gomock.Any(), gomock.Eq(shortened)).
		Return(url, nil).
		Times(1)

	ss.mocks.logger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	res, err := New(ss.mocks.repo).GetURL(ss.ctx, shortened)
	ss.Require().NoError(err)
	ss.Require().Equal(url, res)
}
