package usecase

import (
	"context"
	usecase "github.com/btc-price/internal/usecase/providerChain"
	"github.com/btc-price/pkg/btcpricelb"
)

type Provider struct {
	provider ProviderClient
	next     usecase.Chain
}

func NewProvider(provider ProviderClient, next usecase.Chain) *Provider {
	return &Provider{provider: provider, next: next}
}

type ProviderClient interface {
	GetRateRequest(ctx context.Context, bCurr string, qCurr string) (btcpricelb.Rate, error)
}

func (p *Provider) GetRate(ctx context.Context, bCurr string, qCurr string) (btcpricelb.Rate, error) {
	rate, err := p.provider.GetRateRequest(ctx, bCurr, qCurr)
	if err != nil {
		if p.next == nil {
			return 0, err
		}
		rate, err = p.next.GetRate(ctx, bCurr, qCurr)
	}

	return rate, nil
}
