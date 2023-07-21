package usecase

import (
	"context"
	"fmt"
	"github.com/btc-price/pkg/btcpricelb"
)

type (
	RateService struct {
		rateProvider RateProvider
	}

	RateProvider interface {
		GetCurrencyRate(ctx context.Context, bCurr, qCurr string) (btcpricelb.Rate, error)
	}
)

func NewRateService(rp RateProvider) *RateService {
	return &RateService{rateProvider: rp}
}

func (s *RateService) HandleRate(ctx context.Context, bCurr string, qCurr string) (btcpricelb.RateResponse, error) {
	rate, err := s.rateProvider.GetCurrencyRate(ctx, bCurr, qCurr)
	if err != nil {
		return btcpricelb.RateResponse{}, fmt.Errorf("get rate: %w", err)
	}

	return btcpricelb.RateResponse{Rate: float64(rate)}, nil
}
