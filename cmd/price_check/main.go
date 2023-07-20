package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/btc-price/internal/binanceclient"
	"github.com/btc-price/internal/coingeckoclient"
	"github.com/btc-price/internal/mailing"
	"github.com/btc-price/internal/mailing/emailcreator"
	"github.com/btc-price/internal/mailing/emailsender"
	"github.com/btc-price/internal/rate"
	"github.com/btc-price/internal/rateprovider"
	"github.com/btc-price/internal/subscription"
	"github.com/btc-price/internal/subscription/storage"

	"github.com/btc-price/pkg/coinconverter"

	"github.com/btc-price/cmd/price_check/handler"
	"github.com/caarlos0/env/v6"

	"go.uber.org/zap"
)

func main() {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		panic(fmt.Errorf("read config: %w", err))
	}

	ccx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, stop := signal.NotifyContext(ccx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger, _ := zap.NewProduction()
	defer logger.Sync() //nolint:errcheck

	httpCl := &http.Client{}
	cnv := coinconverter.NewConverter()

	coingeckoCl := coingeckoclient.NewClient(httpCl, cfg.Coingecko.RatePath, cnv)
	binanceCl := binanceclient.NewClient(httpCl, cfg.Binance.RatePath, cnv)

	binancePrvdr := rateprovider.NewLoggingProvider(rateprovider.NewProvider(binanceCl, nil), logger)
	cngkPrvdr := rateprovider.NewLoggingProvider(rateprovider.NewProvider(coingeckoCl, binancePrvdr), logger)

	rateSrv := rate.NewRateService(cngkPrvdr)

	emailStorage := storage.NewStorage(cfg.EmailStorage.Path)

	sbscrSrv := subscription.NewSubscriptionService(emailStorage)

	ec := emailcreator.NewEmailCreator(cfg.EmailText.Text)

	mailingSrv := mailing.NewMailingService(
		emailsender.NewSender(),
		emailStorage,
		cngkPrvdr,
		ec)

	btcPriceHndlr := handler.NewBtcPrice(
		rateSrv,
		sbscrSrv,
		mailingSrv,
		logger)

	router := handler.MakeRouter(ctx, btcPriceHndlr)

	httpServer := &http.Server{
		Addr:           cfg.Port,
		Handler:        router,
		ReadTimeout:    cfg.ServerTimeout,
		WriteTimeout:   cfg.ServerTimeout,
		IdleTimeout:    cfg.ServerTimeout,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}

	errCh := make(chan error)

	_, err := os.Stat(cfg.EmailStorage.Path)
	if err != nil {
		file, errF := os.Create(cfg.EmailStorage.Path)
		defer file.Close()
		if errF != nil {
			logger.Fatal("create file", zap.Error(errF))
		}
	}

	go func() {
		logger.Info("listen and serve", zap.String("address", cfg.Port))

		if err := httpServer.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	shutdown := func() {
		stop()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second) //nolint:gomnd
		defer cancel()

		if err := httpServer.Shutdown(ctxShutdown); err != nil {
			logger.Error("http server: shutdown", zap.Error(err))

			return
		}

		logger.Info("service shutdown: graceful!")
	}

	select {
	case err := <-errCh:
		logger.Error("shutdown catch error", zap.Error(err))
		shutdown()
	case <-ctx.Done():
		logger.Info("shutdown context done")
		shutdown()
	}
}
