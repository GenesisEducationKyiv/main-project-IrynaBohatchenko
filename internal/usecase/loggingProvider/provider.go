package usecase

import (
	"context"
	"fmt"
	usecase "github.com/btc-price/internal/usecase/providerChain"
	"github.com/btc-price/pkg/btcpricelb"
	"go.uber.org/zap"
	"reflect"
)

type LoggingProvider struct {
	provider usecase.Chain
	logger   *zap.Logger
}

func NewLoggingProvider(provider usecase.Chain, logger *zap.Logger) *LoggingProvider {
	return &LoggingProvider{provider: provider, logger: logger}
}

func (l *LoggingProvider) GetRate(ctx context.Context, bCurr string, qCurr string) (btcpricelb.Rate, error) {
	rate, err := l.provider.GetRate(ctx, bCurr, qCurr)
	if err != nil {
		return 0, err
	}

	l.logger.Info(fmt.Sprintf("%s rate response: %f", l.getEProviderName(), rate))

	return rate, nil
}

func (l *LoggingProvider) getEProviderName() string {
	return reflect.TypeOf(l.provider).Elem().Name()
}
