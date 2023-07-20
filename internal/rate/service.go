package rate

import (
	"context"
	"fmt"
	"github.com/btc-price/internal/rateprovider"
	"github.com/btc-price/pkg/btcpricelb"
)

type (
	RateService struct {
		rateProvider rateprovider.RateProvider
	}
)

func NewRateService(rp rateprovider.RateProvider) *RateService {
	return &RateService{rateProvider: rp}
}

func (s *RateService) HandleRate(ctx context.Context, bCurr string, qCurr string) (btcpricelb.RateResponse, error) {
	rate, err := s.rateProvider.GetCurrencyRate(ctx, bCurr, qCurr)
	if err != nil {
		return btcpricelb.RateResponse{}, fmt.Errorf("get rate: %w", err)
	}

	return btcpricelb.RateResponse{Rate: float64(rate)}, nil
}
