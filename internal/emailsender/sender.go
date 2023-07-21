package emailsender

import (
	"context"
	"github.com/btc-price/internal/models"
)

type Sender struct {
}

func NewSender() *Sender {
	return &Sender{}
}

func (s *Sender) SendEmails(ctx context.Context, emailsList []models.Email, text string) error {
	// interaction with third-party email sending service
	// for _, _ = range emailsList {
	// }
	return nil
}
