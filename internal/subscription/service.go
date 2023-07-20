package subscription

import (
	"context"
	"fmt"
	"github.com/btc-price/internal/models"
	"github.com/btc-price/internal/subscription/storageerrors"
	"net/mail"
)

type SubscriptionService struct {
	emailStorage SubscriptionStorage
}

func NewSubscriptionService(emailStorage SubscriptionStorage) *SubscriptionService {
	return &SubscriptionService{emailStorage: emailStorage}
}

type SubscriptionStorage interface {
	AddUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, user *models.User) bool
}

func (s *SubscriptionService) HandleSubscribe(ctx context.Context, email string) error {
	if !s.validateEmail(models.Email(email)) {
		return storageerrors.ErrInvalidEmail
	}

	if s.emailStorage.GetUser(ctx, models.NewUser(models.Email(email))) {
		return storageerrors.ErrEmailExists
	}

	if err := s.emailStorage.AddUser(ctx, models.NewUser(models.Email(email))); err != nil {
		return fmt.Errorf("error adding email %w", err)
	}

	return nil
}

func (s *SubscriptionService) validateEmail(email models.Email) bool {
	_, err := mail.ParseAddress(string(email))
	return err == nil
}
