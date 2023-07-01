package mailing

import (
	"context"
	"fmt"
	"github.com/btc-price/pkg/btcpricelb"
)

type Service struct {
	sender  Sender
	storage EmailStorage
}

func NewService(sender Sender, storage EmailStorage) *Service {
	return &Service{sender: sender, storage: storage}
}

type Sender interface {
	SendEmails(ctx context.Context, emailsList []btcpricelb.Email) error
}

type EmailStorage interface {
	ReadOneEmail(ctx context.Context, email btcpricelb.Email) bool
	ReadAllEmails(ctx context.Context) ([]btcpricelb.Email, error)
}

func (s *Service) HandleSendEmails(ctx context.Context) error {
	list, err := s.storage.ReadAllEmails(ctx)
	if err != nil {
		return fmt.Errorf("read emails: %w", err)
	}

	if err = s.sender.SendEmails(ctx, list); err != nil {
		return fmt.Errorf("send emails: %w", err)
	}

	return nil
}
