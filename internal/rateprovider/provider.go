package rateprovider

import (
	"context"
	"github.com/btc-price/pkg/btcpricelb"
)

type Provider struct {
	provider ProviderClient
	next     RateProvider
}

func NewProvider(provider ProviderClient, next RateProvider) *Provider {
	return &Provider{provider: provider, next: next}
}

type ProviderClient interface {
	GetRateRequest(ctx context.Context, bCurr string, qCurr string) (btcpricelb.Rate, error)
}

type RateProvider interface {
	GetCurrencyRate(ctx context.Context, bCurr, qCurr string) (btcpricelb.Rate, error)
}

func (p *Provider) GetCurrencyRate(ctx context.Context, bCurr string, qCurr string) (btcpricelb.Rate, error) {
	rate, err := p.provider.GetRateRequest(ctx, bCurr, qCurr)
	if err != nil {
		if p.next == nil {
			return 0, err
		}
		rate, err = p.next.GetCurrencyRate(ctx, bCurr, qCurr)
	}

	return rate, nil
}
