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

	"github.com/btc-price/cmd/price_check/handler"
	"github.com/btc-price/internal/btcpriceservice"
	"github.com/btc-price/internal/coingeckoclient"
	"github.com/btc-price/internal/emailsender"
	"github.com/btc-price/internal/emailstorage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func makeTestRouter(cfg Config) http.Handler {
	logger, _ := zap.NewProduction()
	defer logger.Sync() //nolint:errcheck

	btcPriceSrv := btcpriceservice.NewService(
		coingeckoclient.NewClient(cfg.Coingecko.RatePath),
		emailstorage.NewStorage(cfg.EmailStorage.Path),
		emailsender.NewSender())

	btcPriceHndlr := handler.NewBtcPrice(
		btcPriceSrv,
		logger)

	ctx := context.Background()

	return handler.MakeRouter(ctx, btcPriceHndlr)
}

func Test_handle_rate(t *testing.T) {
	t.Parallel()
	cfg := Config{Coingecko: Coingecko{
		RatePath: "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=uah",
	}}

	tempSrv := httptest.NewServer(makeTestRouter(cfg))
	t.Cleanup(func() { tempSrv.Close() })

	t.Run("check getting rate", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/rate", tempSrv.URL), nil)

		if err != nil {
			t.Fatalf("new rate request %v", err)
		}

		client := &http.Client{}

		resp, err := client.Do(req.WithContext(context.Background()))

		if err != nil {
			t.Fatalf("getting rate %v", err)
		}

		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func Test_handle_subscribe(t *testing.T) {
	t.Parallel()

	const filePath = "./emails_test.txt"
	const email = "test_email@gmail.com"

	cfg := Config{
		EmailStorage: FileStorage{Path: filePath},
	}

	tempSrv := httptest.NewServer(makeTestRouter(cfg))
	t.Cleanup(func() { tempSrv.Close() })

	_, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("create file %s", err)
	}
	defer os.Remove(filePath)

	t.Run("check subscribing", func(t *testing.T) {
		t.Parallel()

		data := url.Values{"email": {email}}
		req, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s/api/subscribe", tempSrv.URL),
			bytes.NewBufferString(data.Encode()))

		if err != nil {
			t.Fatalf("new rate request %v", err)
		}

		client := &http.Client{}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req.WithContext(context.Background()))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
