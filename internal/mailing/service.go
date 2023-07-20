package mailing

import (
	"context"
	"fmt"
	"github.com/btc-price/internal/models"
	"github.com/btc-price/internal/rateprovider"
	"github.com/btc-price/pkg/btcpricelb"
)

type MailingService struct {
	sender       Sender
	storage      EmailStorage
	rateProvider rateprovider.RateProvider
	emailCreator EmailCreator
}

func NewMailingService(sender Sender, storage EmailStorage, rp rateprovider.RateProvider, ec EmailCreator) *MailingService {
	return &MailingService{
		sender:       sender,
		storage:      storage,
		rateProvider: rp,
		emailCreator: ec,
	}
}

type Sender interface {
	SendEmails(ctx context.Context, emailsList []models.Email, text string) error
}

type EmailStorage interface {
	GetUsersList(ctx context.Context) ([]*models.User, error)
}

type EmailCreator interface {
	GenerateEmail(ctx context.Context, rate btcpricelb.Rate) string
}

func (s *MailingService) HandleSendEmails(ctx context.Context, bCurr, qCurr string) error {
	list, err := s.storage.GetUsersList(ctx)
	if err != nil {
		return fmt.Errorf("read emails: %w", err)
	}

	rate, err := s.rateProvider.GetCurrencyRate(ctx, bCurr, qCurr)
	if err != nil {

	}
	emailText := s.emailCreator.GenerateEmail(ctx, rate)

	if err = s.sender.SendEmails(ctx, s.getEmailsList(list), emailText); err != nil {
		return fmt.Errorf("send emails: %w", err)
	}

	return nil
}

func (s *MailingService) getEmailsList(usersList []*models.User) []models.Email {
	emailsList := []models.Email{}
	for _, user := range usersList {
		emailsList = append(emailsList, user.Email)
	}

	return emailsList
}
