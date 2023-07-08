package usecase

import (
	"context"
	"fmt"
	"github.com/btc-price/pkg/btcpricelb"
)

type ProvidersChain struct {
	firstElement Chain
}

func NewProvidersChain(firstElement Chain) *ProvidersChain {
	return &ProvidersChain{firstElement: firstElement}
}

type Chain interface {
	GetRate(ctx context.Context, bCurr string, qCurr string) (btcpricelb.Rate, error)
}

func (p *ProvidersChain) GetCurrencyRate(ctx context.Context, bCurr string, qCurr string) (btcpricelb.Rate, error) {
	rate, err := p.firstElement.GetRate(ctx, bCurr, qCurr)
	if err != nil {
		return 0, fmt.Errorf("get rate error: %w", err)
	}

	return rate, nil
}
