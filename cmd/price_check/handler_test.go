// E2E test of endpoints
package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/btc-price/internal/binanceclient"
	emailCreatorusecase "github.com/btc-price/internal/usecase/emailCreator"
	logPrUsecase "github.com/btc-price/internal/usecase/loggingProvider"
	mailingUsecase "github.com/btc-price/internal/usecase/mailing"
	prUsecase "github.com/btc-price/internal/usecase/provider"
	prChainUsecase "github.com/btc-price/internal/usecase/providerChain"
	rateUsecase "github.com/btc-price/internal/usecase/rate"
	subscriptionUsecase "github.com/btc-price/internal/usecase/subscription"
	"github.com/btc-price/pkg/coinconverter"

	"github.com/btc-price/cmd/price_check/handler"
	"github.com/btc-price/internal/coingeckoclient"
	"github.com/btc-price/internal/emailsender"
	"github.com/btc-price/internal/storage"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type HandlerSuite struct {
	suite.Suite
	Server   *httptest.Server
	Ctx      context.Context
	FilePath string
}

func makeTestRouter(cfg Config) http.Handler {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() //nolint:errcheck

	httpCl := &http.Client{}
	cnv := coinconverter.NewConverter()

	coingeckoCl := coingeckoclient.NewClient(httpCl, cfg.Coingecko.RatePath, cnv)
	binanceCl := binanceclient.NewClient(httpCl, cfg.Binance.RatePath, cnv)

	binancePrvdr := logPrUsecase.NewLoggingProvider(prUsecase.NewProvider(binanceCl, nil), logger)
	cngkPrvdr := logPrUsecase.NewLoggingProvider(prUsecase.NewProvider(coingeckoCl, binancePrvdr), logger)

	providersChain := prChainUsecase.NewProvidersChain(cngkPrvdr)

	rateSrv := rateUsecase.NewRateService(providersChain)

	emailStorage := storage.NewStorage(cfg.EmailStorage.Path)

	sbscrSrv := subscriptionUsecase.NewSubscriptionService(emailStorage)

	ec := emailCreatorusecase.NewEmailCreator(cfg.EmailText.Text)

	mailingSrv := mailingUsecase.NewMailingService(
		emailsender.NewSender(),
		emailStorage,
		providersChain,
		ec)

	btcPriceHndlr := handler.NewBtcPrice(
		rateSrv,
		sbscrSrv,
		mailingSrv,
		logger)

	ctx := context.Background()

	return handler.MakeRouter(ctx, btcPriceHndlr)
}

func (s *HandlerSuite) SetupSuite() {
	s.FilePath = "./emails_test.txt"
	cfg := Config{
		Coingecko: Coingecko{
			RatePath: "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=uah",
		},
		EmailStorage: FileStorage{Path: s.FilePath},
	}
	s.Server = httptest.NewServer(makeTestRouter(cfg))
	s.Ctx = context.Background()
}

func (s *HandlerSuite) TearDownSuite() {
	s.Server.Close()
}

func (s *HandlerSuite) TestHandleRate() {
	s.Run("check getting rate", func() {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/rate", s.Server.URL), nil)

		if err != nil {
			s.FailNowf("", "new rate request %s", err.Error())
		}

		client := &http.Client{}

		resp, err := client.Do(req.WithContext(context.Background()))

		if err != nil {
			s.FailNowf("", "getting rate %s", err.Error())
		}

		defer resp.Body.Close()

		s.Equal(http.StatusOK, resp.StatusCode)
	})
}

func (s *HandlerSuite) TestHandleSubscribe() {
	_, err := os.Create(s.FilePath)
	if err != nil {
		s.FailNowf("", "create file %s", err.Error())
	}
	defer os.Remove(s.FilePath)

	testCases := []struct {
		name  string
		email string
		want  int
	}{
		{
			name:  "successful subscribing",
			email: "test_email@gmail.com",
			want:  http.StatusOK,
		},
		{
			name:  "fail subscribing with invalid email",
			email: "test_email",
			want:  http.StatusBadRequest,
		},
	}

	for _, tt := range testCases {
		s.Run(tt.name, func() {
			data := url.Values{"email": {tt.email}}
			req, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("%s/api/subscribe", s.Server.URL),
				bytes.NewBufferString(data.Encode()))

			if err != nil {
				s.FailNowf("", "new subscribe request %s", err.Error())
			}

			client := &http.Client{}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := client.Do(req.WithContext(context.Background()))
			if err != nil {
				s.FailNowf("", "Expected no error, got %s", err.Error())
			}

			defer resp.Body.Close()

			s.Equal(tt.want, resp.StatusCode)
		})
	}
}

func TestHandler(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HandlerSuite))
}
