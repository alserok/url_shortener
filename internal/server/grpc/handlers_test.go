package grpc

import (
	"context"
	"github.com/alserok/url_shortener/internal/cache"
	"github.com/alserok/url_shortener/internal/service"
	"github.com/alserok/url_shortener/internal/service/models"
	"github.com/alserok/url_shortener/internal/utils"
	"github.com/alserok/url_shortener/pkg/logger"
	"github.com/alserok/url_shortener/pkg/proto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestGRPCHandlersSuite(t *testing.T) {
	suite.Run(t, new(GRPCHandlersSuite))
}

type GRPCHandlersSuite struct {
	suite.Suite

	handler handler
	ctx     context.Context

	mocks struct {
		ctrl *gomock.Controller

		cache  *cache.MockCache
		srvc   *service.MockService
		logger *logger.MockLogger
	}
}

func (ghs *GRPCHandlersSuite) SetupTest() {
	ghs.mocks.ctrl = gomock.NewController(ghs.T())
	ghs.mocks.srvc = service.NewMockService(ghs.mocks.ctrl)
	ghs.mocks.cache = cache.NewMockCache(ghs.mocks.ctrl)
	ghs.mocks.logger = logger.NewMockLogger(ghs.mocks.ctrl)

	ghs.handler = handler{srvc: ghs.mocks.srvc, cache: ghs.mocks.cache}
	ghs.ctx = logger.WrapLogger(context.Background(), ghs.mocks.logger)
}

func (ghs *GRPCHandlersSuite) TeardownRest() {
	ghs.mocks.ctrl.Finish()
}

func (ghs *GRPCHandlersSuite) TestShortenAndSaveURL() {
	url := models.URL{OriginURL: "url"}
	shortenedURL := "shortened"

	ghs.mocks.srvc.EXPECT().
		ShortenAndSaveURL(gomock.Any(), gomock.Eq(url.OriginURL)).
		Return("shortened", nil).
		Times(1)

	ghs.mocks.logger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	res, err := ghs.handler.ShortenAndSaveURL(ghs.ctx, &proto.URL{OriginUrl: url.OriginURL})
	ghs.Require().NoError(err)
	ghs.Require().Equal(shortenedURL, res.ShortenedUrl)
}

func (ghs *GRPCHandlersSuite) TestGetURL() {
	shortenedURL := "shortened"
	url := "url"

	ghs.mocks.srvc.EXPECT().
		GetURL(gomock.Any(), gomock.Eq(shortenedURL)).
		Return(url, nil).
		Times(1)

	ghs.mocks.cache.EXPECT().
		Get(gomock.Any(), gomock.Eq(shortenedURL)).
		Return("", utils.NewError("not found", utils.NotFoundErr)).
		Times(1)

	ghs.mocks.cache.EXPECT().
		Set(gomock.Any(), gomock.Eq(shortenedURL), gomock.Eq(url)).
		Return(nil).
		Times(1)

	ghs.mocks.logger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	res, err := ghs.handler.GetURL(ghs.ctx, &proto.ShortenedURL{ShortenedUrl: shortenedURL})
	ghs.Require().NoError(err)
	ghs.Require().Equal(url, res.OriginUrl)
}
