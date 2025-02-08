package service

import (
	"context"
	"github.com/alserok/url_shortener/internal/db"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite

	mocks struct {
		ctrl *gomock.Controller

		repo *db.MockRepository
	}
}

func (ss *ServiceSuite) SetupTest() {
	ss.mocks.ctrl = gomock.NewController(ss.T())
	ss.mocks.repo = db.NewMockRepository(ss.mocks.ctrl)
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

	res, err := New(ss.mocks.repo).ShortenAndSaveURL(context.Background(), url)
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

	res, err := New(ss.mocks.repo).GetURL(context.Background(), shortened)
	ss.Require().NoError(err)
	ss.Require().Equal(url, res)
}
