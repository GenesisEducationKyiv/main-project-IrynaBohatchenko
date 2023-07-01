package subscription

import (
	"context"
	"fmt"
	"github.com/btc-price/pkg/btcpricelb"
)

type Service struct {
	emailStorage SubscriptionStorage
}

func NewService(emailStorage SubscriptionStorage) *Service {
	return &Service{emailStorage: emailStorage}
}

type SubscriptionStorage interface {
	AddEmail(ctx context.Context, email btcpricelb.Email) error
	ReadOneEmail(ctx context.Context, email btcpricelb.Email) bool
}

func (s *Service) HandleSubscribe(ctx context.Context, email string) error {
	if err := s.emailStorage.AddEmail(ctx, btcpricelb.Email(email)); err != nil {
		return fmt.Errorf("error adding email %w", err)
	}

	return nil
}
