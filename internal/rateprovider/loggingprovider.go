package rateprovider

import (
	"context"
	"fmt"
	"github.com/btc-price/pkg/btcpricelb"
	"go.uber.org/zap"
	"reflect"
)

type LoggingProvider struct {
	provider RateProvider
	logger   *zap.Logger
}

func NewLoggingProvider(provider RateProvider, logger *zap.Logger) *LoggingProvider {
	return &LoggingProvider{provider: provider, logger: logger}
}

func (l *LoggingProvider) GetCurrencyRate(ctx context.Context, bCurr string, qCurr string) (btcpricelb.Rate, error) {
	rate, err := l.provider.GetCurrencyRate(ctx, bCurr, qCurr)
	if err != nil {
		return 0, err
	}

	l.logger.Info(fmt.Sprintf("%s rate response: %f", l.getEProviderName(), rate))

	return rate, nil
}

func (l *LoggingProvider) getEProviderName() string {
	return reflect.TypeOf(l.provider).Elem().Name()
}
