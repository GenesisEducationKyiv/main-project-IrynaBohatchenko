package rate

import (
	"context"
	"fmt"
	"github.com/btc-price/pkg/btcpricelb"
)

type (
	Service struct {
		coingecko CoingeckoClient
	}

	CoingeckoClient interface {
		GetRate(ctx context.Context, marketCurr string, baseCurr string) (btcpricelb.CoingeckoRate, error)
	}
)

func NewService(coingecko CoingeckoClient) *Service {
	return &Service{coingecko: coingecko}
}

func (s *Service) HandleRate(ctx context.Context, marketCurr string, baseCurr string) (btcpricelb.RateResponse, error) {
	rate, err := s.coingecko.GetRate(ctx, marketCurr, baseCurr)
	if err != nil {
		return btcpricelb.RateResponse{}, fmt.Errorf("get rate: %w", err)
	}

	return btcpricelb.RateResponse{Rate: float64(rate)}, nil
}
