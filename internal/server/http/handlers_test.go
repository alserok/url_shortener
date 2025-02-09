package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/alserok/url_shortener/internal/cache"
	"github.com/alserok/url_shortener/internal/service"
	"github.com/alserok/url_shortener/internal/service/models"
	"github.com/alserok/url_shortener/internal/utils"
	"github.com/alserok/url_shortener/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPHandlersSuite(t *testing.T) {
	suite.Run(t, new(HTTPHandlersSuite))
}

type HTTPHandlersSuite struct {
	suite.Suite

	handler handler
	ctx     context.Context

	mocks struct {
		ctrl *gomock.Controller

		srvc   *service.MockService
		cache  *cache.MockCache
		logger *logger.MockLogger
	}
}

func (hhs *HTTPHandlersSuite) SetupTest() {
	hhs.mocks.ctrl = gomock.NewController(hhs.T())
	hhs.mocks.srvc = service.NewMockService(hhs.mocks.ctrl)
	hhs.mocks.cache = cache.NewMockCache(hhs.mocks.ctrl)
	hhs.mocks.logger = logger.NewMockLogger(hhs.mocks.ctrl)

	hhs.handler = handler{srvc: hhs.mocks.srvc, cache: hhs.mocks.cache}
	hhs.ctx = logger.WrapLogger(context.Background(), hhs.mocks.logger)
}

func (hhs *HTTPHandlersSuite) TeardownRest() {
	hhs.mocks.ctrl.Finish()
}

func (hhs *HTTPHandlersSuite) TestShortenAndSaveURL() {
	url := models.URL{OriginURL: "url"}
	shortenedURL := "shortened"

	hhs.mocks.srvc.EXPECT().
		ShortenAndSaveURL(gomock.Any(), gomock.Eq(url.OriginURL)).
		Return("shortened", nil).
		Times(1)

	hhs.mocks.logger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	b, err := json.Marshal(url)
	hhs.Require().NoError(err)
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(b))
	r = r.WithContext(hhs.ctx)
	w := httptest.NewRecorder()

	hhs.Require().NoError(hhs.handler.ShortenAndSaveURL(w, r))

	var res map[string]any
	hhs.Require().NoError(json.NewDecoder(w.Body).Decode(&res))
	hhs.Require().Equal(shortenedURL, res["shortenedURL"])
}

func (hhs *HTTPHandlersSuite) TestGetURL() {
	shortenedURL := "shortened"
	url := "url"

	hhs.mocks.srvc.EXPECT().
		GetURL(gomock.Any(), gomock.Eq(shortenedURL)).
		Return(url, nil).
		Times(1)

	hhs.mocks.logger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	hhs.mocks.cache.EXPECT().
		Get(gomock.Any(), gomock.Eq(shortenedURL)).
		Return("", utils.NewError("not found", utils.NotFoundErr)).
		Times(1)

	hhs.mocks.cache.EXPECT().
		Set(gomock.Any(), gomock.Eq(shortenedURL), gomock.Eq(url)).
		Return(nil).
		Times(1)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r = r.WithContext(hhs.ctx)
	r.SetPathValue("shortenedURL", shortenedURL)
	w := httptest.NewRecorder()

	hhs.Require().NoError(hhs.handler.GetURL(w, r))

	var res map[string]any
	hhs.Require().NoError(json.NewDecoder(w.Body).Decode(&res))
	hhs.Require().Equal(url, res["originURL"])
}
